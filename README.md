# Go-Do-This Task Manager

A small command-line to-do list written in Go. Tasks are stored in a local JSON file, so your list persists between runs.

## Features

- Add, edit, delete, and list tasks
- Mark tasks as complete or reset them to incomplete
- Tasks persist in tasks.json in the directory you run the program from
- Output is formatted as aligned tables
- Uses only the Go standard library (no dependencies)

## Project Structure

main.go  - Entry point. Parses command-line arguments, dispatches commands, and prints the usage table
tasks.go - Defines the Task and TaskList types, the operations on them (add, complete, reset, edit, delete, print), and JSON load/save

> The file names above are suggestions. Rename them to match your actual files.

## Getting Started

To use the list use terminal commands of the format:

go run . <command>

## Usage

Running the program with no arguments (go run .) prints the command table:

Command: Action
print: show all tasks
add <name>: add a task
done <id>: mark a task complete
reset <id>: mark a task incomplete
delete <id>: remove a task
edit <id> <new name>: edit a task

Task names can contain spaces; no quotes are needed.

### Examples

go run . add Buy groceries
Added: | 1 | Buy groceries |

go run . add Finish report
Added: | 2 | Finish report |

go run . done 1
Completed: | 1 | Buy groceries |

go run . edit 2 Finish quarterly report
Edited: | 2 | Finish quarterly report |

go run . print


The print command shows a table of the form:


| ID | Task Name                | Completion |
| 1  | Buy groceries            |    [--]    |
| 2  | Finish quarterly report  |    [  ]    |


A [--] mark means the task is complete.

## How It Works

1. Startup: main calls run with the command-line arguments (excluding the program name). Any error is printed to stderr and the program exits with status 1.
2. Loading: Load reads tasks.json. If the file doesn't exist, an empty list is used. Because the JSON stores only the tasks, the highest existing ID is recomputed on load so new tasks always receive a unique ID.
3. Command dispatch: A switch on the first argument selects the operation. IDs are parsed by parseID, which validates that an ID was supplied and is a number.
4. Saving: After any command that changes the list (everything except print), Save writes the tasks back to tasks.json as indented JSON.

### Data Format

Tasks are stored as a JSON array:

json

{
 "id": 1,
 "name": "Buy water",
 "done": true
}

IDs are never reused within a session: deleting a task does not shift the IDs of the others. Note that if you delete the task with the highest ID and then reload, that ID can be reused, since the next ID is derived from the highest remaining one.

## Error Handling

- Unknown commands, missing arguments, non-numeric IDs, and IDs that don't match any task produce a message such as error: task 5 not found and exit with a non-zero status.
- A missing tasks.json is not an error; it is created on the first save.
