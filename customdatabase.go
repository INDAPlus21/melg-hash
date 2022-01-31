// go build -o C:\Users\Marcus\Documents\Github\melg-hash
// .\customdatabase.exe

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var database map[string]string = make(map[string]string)

func main() {
	fmt.Println("Loading database...")
	load()
	get_input()
}

func get_input() {
	for {
		fmt.Println("\nEnter a command")

		// Read input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\r\n")
		values := strings.Split(input, " ")

		if len(values) == 0 {
			fmt.Println("Invalid command")
		} else {
			// Different commands
			if values[0] == "ADD" {
				if len(values) != 3 {
					fmt.Println("Invalid input, format is: ADD <KEY> <VALUE>")
				} else {
					database[values[1]] = values[2]
					fmt.Println("Sucessfully added entry with key", values[1], "and value", values[2])
				}
			} else if values[0] == "REMOVE" {
				if len(values) != 2 {
					fmt.Println("Invalid input, format is: REMOVE <KEY>")
				} else {
					_, ok := database[values[1]]
					delete(database, values[1])
					if ok {
						fmt.Println("Sucessfully removed entry with key", values[1])
					} else {
						fmt.Println("Did not find entry with key", values[1])
					}
				}
			} else if values[0] == "GET" {
				if len(values) != 2 {
					fmt.Println("Invalid input, format is: GET <KEY>")
				} else {
					value, ok := database[values[1]]
					if ok {
						fmt.Println("Key", values[1], "was found in database with value", value)
					} else {
						fmt.Println("Key", values[1], "was not found in database")
					}
				}
			} else if values[0] == "LIST" {
				fmt.Println("Database key and values:")
				for key, value := range database {
					fmt.Println(key + ":" + value)
				}
			} else {
				fmt.Println("Invalid command")
			}
		}

		save()
	}
}

func load() {
	// Check file
	directory, _ := os.Getwd()
	file_location := directory + "\\data.data"
	file, err := os.Open(file_location)

	if os.IsNotExist(err) {
		fmt.Println("No file found to load data from")
		return
	} else if err != nil {
		panic(err)
	}

	// Load data line for line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ":")
		if len(values) != 2 {
			fmt.Println("Line is not formatted correctly:", values)
			continue
		}
		database[values[0]] = values[1]
	}

	file.Close()
}

func save() {
	// Create (override) file
	directory, _ := os.Getwd()
	file_location := directory + "\\data.data"
	file, err := os.Create(file_location)

	if err != nil {
		fmt.Println("Error when creating file")
	}

	// Write all database data
	writer := bufio.NewWriter(file)
	for key, value := range database {
		writer.WriteString(key + ":" + value + "\n")
	}

	writer.Flush()
	file.Close()
}
