package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	file = "tasks.csv"
	help = "help.txt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No command provided. Use 'help' for usage information.")
	}
	command := os.Args[1]
	switch command {
	case "list":
		err := listTasks()
		if err != nil {
			log.Fatal(err.Error())
		}
	case "help":
		text, err := readFile(help)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(text))
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("No ID provided. Use 'help' for usage information.")
		}
		providedId := os.Args[2]
		err := deleteTask(providedId)
		if err != nil {
			log.Fatal(err.Error())
		}
	case "complete":
		if len(os.Args) < 3 {
			log.Fatal("No ID provided. Use 'help' for usage information.")
		}
		providedId := os.Args[2]
		err := completeTask(providedId)
		if err != nil {
			log.Fatal(err.Error())
		}
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("No description provided. Use 'help' for usage information.")
		}
		description := os.Args[2]
		err := addTask(description)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func listTasks() error {
	tasks, err := readFile(file)
	if err != nil {
		return err
	}
	fmt.Printf("\n%-10s %-22s %-15s %-15s\n", "ID", "Description", "Completed", "Created At")
	fmt.Println(strings.Repeat("—", 69))
	rows := strings.Split(string(tasks), "\n")
	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			rowSlice := strings.Split(row, ",")
			if len(rowSlice) == 4 {
				fmt.Printf("%-10s %-22s %-15s %-15s\n", rowSlice[0], rowSlice[1], rowSlice[2], rowSlice[3])
			}
		}
	}
	fmt.Println()
	return nil
}

func deleteTask(id string) error {
	tasks, err := readFile(file)
	if err != nil {
		return err
	}
	rows := strings.Split(string(tasks), "\n")
	var updatedRows []string
	found := false
	for _, row := range rows {
		if strings.TrimSpace(row) != "" && id != strings.Split(row, ",")[0] {
			updatedRows = append(updatedRows, row)
		} else if strings.TrimSpace(row) != "" && id == strings.Split(row, ",")[0] {
			found = true
		}
	}
	err = os.WriteFile(file, []byte(strings.Join(updatedRows, "\n")), 0644)
	if err != nil {
		return err
	}
	if found {
		fmt.Printf("Task with ID %s deleted.\n", id)
	} else {
		fmt.Printf("Task with ID %s was not found.\n", id)
	}
	return nil
}

func completeTask(id string) error {
	tasks, err := readFile(file)
	if err != nil {
		return err
	}
	rows := strings.Split(string(tasks), "\n")
	var updatedRows []string
	found := false
	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			rowSlice := strings.Split(row, ",")
			if len(rowSlice) == 4 && rowSlice[0] == id {
				if rowSlice[2] == "true" {
					rowSlice[2] = "false"
				} else {
					rowSlice[2] = "true"
				}
				found = true
				updatedRows = append(updatedRows, strings.Join(rowSlice, ","))
			} else if len(rowSlice) == 4 {
				updatedRows = append(updatedRows, row)
			}
		}
	}
	err = os.WriteFile(file, []byte(strings.Join(updatedRows, "\n")), 0644)
	if err != nil {
		return err
	}
	if found {
		fmt.Printf("Task with ID %s toggled completion status.\n", id)
	} else {
		fmt.Printf("Task with ID %s was not found.\n", id)
	}
	return nil
}

func addTask(description string) error {
	tasks, err := readFile(file)
	if err != nil {
		return err
	}
	rows := strings.Split(string(tasks), "\n")
	nextId, err := getNextId()
	if err != nil {
		return err
	}
	newTask := []string{strconv.Itoa(nextId), description, "false", time.Now().Format(time.DateTime)}
	rows = append(rows, strings.Join(newTask, ","))
	err = os.WriteFile(file, []byte(strings.Join(rows, "\n")), 0644)
	if err != nil {
		return err
	}
	err = saveNextId(nextId + 1)
	if err != nil {
		return err
	}
	fmt.Println("Task added successfully.")
	return nil
}

func getNextId() (int, error) {
	data, err := os.ReadFile("nextid.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil // Start from 1 if the file doesn't exist
		}
		return 0, err
	}
	id, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func saveNextId(id int) error {
	return os.WriteFile("nextid.txt", []byte(strconv.Itoa(id)), 0644)
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			os.Create(filename)
			return []byte{}, nil
		}
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
