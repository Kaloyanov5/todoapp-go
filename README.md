# Todo CLI App

A simple command-line todo application written in Go.

## Features

- Add tasks with descriptions
- List all tasks in a table
- Mark tasks as complete/incomplete
- Delete tasks
- Persistent storage using CSV files

## Usage

Build the app:

```sh
go build -o todoapp
```

Run the commands:

```
./todoapp add "Buy groceries"
./todoapp list
./todoapp complete 1
./todoapp delete 1
./todoapp help
```

## File Structure
- main.go — main application logic
- tasks.csv — stores tasks
- nextid.txt — stores the next task ID
- help.txt — help instructions

## Requirments
- Go 1.24 or newer
