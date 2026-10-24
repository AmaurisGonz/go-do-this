package main

import (
	"fmt"     // Used for printing
	"os"      // Used for command-line arguments
	"strconv" // String to number conversion
	"strings" // More methods for strings
)

const dataFile = "tasks.json" // This is the fixed filename of the json file

// Showcases the guide to using the program
func printUsage() {

	// Here is a list of each command and its action, header is row 0
	rows := [][2]string{
		{"Command", "Action"},
		{"print", "show all tasks"},
		{"add <name>", "add a task"},
		{"done <id>", "mark a task complete"},
		{"reset <id>", "mark a task incomplete"},
		{"delete <id>", "remove a task"},
		{"edit <id> <new name>", "edit a task"},
	}

	// Finds the longest text in each column
	max1, max2 := 0, 0

	// Loops through each row
	for _, r := range rows {

		// Compares first value with max
		if len(r[0]) > max1 {
			max1 = len(r[0])
		}

		// Compares second value with max
		if len(r[1]) > max2 {
			max2 = len(r[1])
		}
	}

	// "| " + col1 + " | " + col2 + " |" adds up to max1 + max2 + 7 characters
	line := strings.Repeat("-", max1+max2+7)

	// Prints each row according to the maxes
	fmt.Println(line)
	for i, r := range rows {
		// Inputs max as the character count for each column
		fmt.Printf("| %-*s | %-*s |\n", max1, r[0], max2, r[1])
		if i == 0 {
			fmt.Println(line)
		}
	}
	fmt.Println(line)
}

// Main function
func main() {

	// Run is our function which takes os.Args and returns an error or nil
	if err := run(os.Args[1:]); err != nil {

		// Prints an error statement to os.Stderr then exists the program
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// This function runs the program, allowing for the user to use the to do list
func run(args []string) error {

	// If there are no arguments, show the help statement and return nil
	if len(args) == 0 {
		printUsage()
		return nil
	}

	// Loads task list from the File sends an error if it fails
	list, err := Load(dataFile)
	if err != nil {
		return err
	}

	// Args[0] is the command word, switch picks the matching case
	switch args[0] {

	// Showcases the list using the print function
	case "print":
		list.Print()
		return nil // nothing changed, so skip saving

	// Adds a task using the add function after checking that a name is included in the arguments
	case "add":
		// Checks that the arguement list has at least 2 values
		if len(args) < 2 {
			return fmt.Errorf("usage: add <name>")
		}

		// Adds the second arguement onwards (name)
		addedTask := list.Add(strings.Join(args[1:], " "))
		fmt.Printf("Added: | %d | %s |\n", addedTask.ID, addedTask.Name)

	// Marks a task as done
	case "done":

		// Uses parse ID function to return the ID
		id, err := parseID(args)
		if err != nil {
			return err
		}

		// Marks the ID complete using the function
		completedTask, err := list.Complete(id)
		if err != nil {
			return err
		}
		fmt.Printf("Completed: | %d | %s | \n", completedTask.ID, completedTask.Name)

	// Marks a task as undone
	case "reset":

		// Uses parse ID function to return the ID
		id, err := parseID(args)
		if err != nil {
			return err
		}

		// Marks the ID incomplete using the function
		completedTask, err := list.Reset(id)
		if err != nil {
			return err
		}
		fmt.Printf("Reset: | %d | %s | \n", completedTask.ID, completedTask.Name)

	case "edit":
		// Checks that arguement list has at least 3 values
		if len(args) < 3 {
			return fmt.Errorf("usage: edit <id> <new name>")
		}

		// Uses parse ID function to return the ID
		id, err := parseID(args)
		if err != nil {
			return err
		}

		// args[2:] is everything after the ID; Join turns it into one string,
		newName := strings.Join(args[2:], " ")

		// Changes the name of the task using the function
		editedTask, err := list.Edit(id, newName)
		if err != nil {
			return err
		}
		fmt.Printf("Edited: | %d | %s | \n", editedTask.ID, editedTask.Name)

	// Removes a task from the tasklist
	case "delete":

		// Uses parseID function to return an ID
		id, err := parseID(args)
		if err != nil {
			return err
		}

		// Removes the ID using the function
		deletedTask, err := list.Delete(id)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted: | %d | %s | \n", deletedTask.ID, deletedTask.Name)

	// This is what is returned when there is no case
	default:

		// Throws an error if the command is not used
		return fmt.Errorf("unknown command %q", args[0])
	}

	// Saves the list
	return list.Save(dataFile)
}

// Converts string to integer while checking that there is a number in the users statement and checking that the string is indeed a number
func parseID(args []string) (int, error) {

	// Checks that there are at least 2 statement
	if len(args) < 2 {
		return 0, fmt.Errorf("missing task ID")
	}

	// Checks that the ID statment is an integer
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, fmt.Errorf("invalid task ID %q", args[1])
	}
	return id, nil
}
