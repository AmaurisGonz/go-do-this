package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// This structure defines what a task is in our to do list
type Task struct {

	// This defines each variable (a unique ID, a name, a variable denoting completion of the task)
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"done"`
}

// This structure defines the complete list of all tasks
type TaskList struct {

	// The is an array of tasks and an ID of the ID of the last item in the list
	Tasks     []Task
	currentID int
}

// The add function adds a name to a tasklist
func (l *TaskList) Add(name string) Task {

	// The ID value of the tasklist is pushed once and then the
	l.currentID++

	// A new task with the name parameter is set to variable and then appended to the tasklist
	t := Task{ID: l.currentID, Name: name}
	l.Tasks = append(l.Tasks, t)

	// The task is returned
	return t
}

// The complete function changed the Complete value on a task based on the ID passed through
func (l *TaskList) Complete(id int) (Task, error) {

	// Here we loop through each ID present in the Tasks list
	for i := range l.Tasks {

		// We check if the ID of the current task is equal to the ID passed through, if so then set done to true and return nil
		if l.Tasks[i].ID == id {
			l.Tasks[i].Complete = true
			return l.Tasks[i], nil
		}
	}

	// This returns an error stating the task was not found if the  loop completed without finding a task
	return Task{}, fmt.Errorf("task %d not found", id)
}

// The reset function marks the task as incomplete
func (l *TaskList) Reset(id int) (Task, error) {

	// Here we loop through each ID present in the Tasks list
	for i := range l.Tasks {

		// We check if the ID of the current task is equal to the ID passed through, if so then set done to true and return nil
		if l.Tasks[i].ID == id {
			l.Tasks[i].Complete = false
			return l.Tasks[i], nil
		}
	}

	// This returns an error stating the task was not found if the  loop completed without finding a task
	return Task{}, fmt.Errorf("task %d not found", id)
}

// The edit function changes the name of a tasked based off of its ID
func (l *TaskList) Edit(id int, newName string) (Task, error) {

	// Here we loop through each ID present in the Tasks list
	for i := range l.Tasks {

		// We check if the ID of the current task is equal to the ID passed through, if so then set the name variable to the new name and return nil
		if l.Tasks[i].ID == id {
			l.Tasks[i].Name = newName
			return l.Tasks[i], nil
		}
	}

	// This returns an error stating the task was not found if the  loop completed without finding a task
	return Task{}, fmt.Errorf("task %d not found", id)
}

// The delete function removes a task from the tasklist based on an id passed through
func (l *TaskList) Delete(id int) (Task, error) {

	// Here we loop through each ID present in the Tasks list
	for i, t := range l.Tasks {

		// We check if the ID of the current task is equal to the ID passed through, if so then we append all the values leading up to the index all the values after it, removing the current value
		if t.ID == id {
			l.Tasks = append(l.Tasks[:i], l.Tasks[i+1:]...)
			return t, nil
		}
	}

	// This returns an error stating the task was not found if the  loop completed without finding a task
	return Task{}, fmt.Errorf("task %d not found", id)
}

// This function prints the tasklist
func (l *TaskList) Print() {

	// Checks if there are no tasks in the list, if so then print accordingly
	if len(l.Tasks) == 0 {
		fmt.Println("No tasks yet.")
		return
	}

	// Sets max to header length
	idMax := len("ID")
	nameMax := len("Task Name")

	// Loops to find the max of each column
	for _, t := range l.Tasks {

		// Compares the ID length to the ID max
		if currentID := len(strconv.Itoa(t.ID)); currentID > idMax {
			idMax = currentID
		}

		// Compares the name length to the name max
		if len(t.Name) > nameMax {
			nameMax = len(t.Name)
		}

	}

	// "| " + ID + " | " + Name + " | " + Completion + " |" is maxID + maxName + 20
	line := strings.Repeat("-", idMax+nameMax+20)

	// Prints header
	fmt.Println(line)
	fmt.Printf("| %*s | %-*s | Completion |\n", idMax, "ID", nameMax, "Task Name")
	fmt.Println(line)

	// Iterates through each task and prints the task along with a symbol representing its completion
	for _, t := range l.Tasks {
		mark := "  "
		if t.Complete {
			mark = "--"
		}
		fmt.Printf("| %-*d | %-*s |    [%s]    |\n", idMax, t.ID, nameMax, t.Name, mark)
	}
	fmt.Println(line)
}

// This function loads a json file containing a tasklist
func Load(path string) (*TaskList, error) {

	// Sets a variable to a pointer that represents an empty tasklist
	list := &TaskList{}

	// This reads the file based on its filepath and returns its contents and an error if one arises
	data, err := os.ReadFile(path)

	// This checks if the file exists by seeing if the error is os.ErrNotExist
	if errors.Is(err, os.ErrNotExist) {

		// If the file does not exist then we return the empty list to make a new file
		return list, nil
	}

	// if the error is not nil then we return a nil list and the error as well
	if err != nil {
		return nil, err
	}

	// We use unmarshal to convert json text to go values and stores it in a given location, the if statement runs the unmarshal first and then checks if there is an error, if so then return nil for the list and the error
	if err := json.Unmarshal(data, &list.Tasks); err != nil {
		return nil, err
	}

	// This fixes the currentID value of the tasklist since the json file does not store the count, it iterates through each task and compares it to the currentID, setting the task ID to the currentID if it is higher than it
	for _, t := range list.Tasks {
		if t.ID > list.currentID {
			list.currentID = t.ID
		}
	}
	return list, nil
}

// This saves the modified tasklist to a json file
func (l *TaskList) Save(path string) error {

	// The marshalindent converts the go data to json data returning an error if it fails
	data, err := json.MarshalIndent(l.Tasks, "", "  ")
	if err != nil {
		return err
	}

	// Here we write then json data to the file aswell as giving it file permissions
	return os.WriteFile(path, data, 0o644)
}
