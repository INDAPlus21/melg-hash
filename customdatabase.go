// go build -o C:\Users\Marcus\Documents\Github\melg-hash
// .\customdatabase.exe

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var database = make(map[string]map[string]string)
var default_database = "default"

func main() {
	fmt.Println("Loading database...")
	load()

	// Add default table
	if len(database) == 0 {
		database[default_database] = make(map[string]string)
	}

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
			switch values[0] {
			case "ADD":
				command_add(values)
			case "REMOVE":
				command_remove(values)
			case "GET":
				command_get(values)
			case "LIST":
				command_list(values)
			case "ADDTABLE":
				command_add_table(values)
			case "LISTTABLES":
				command_list_tables()
			default:
				fmt.Println("Invalid command")
			}
		}

		save()
	}
}

func command_add(values []string) {
	if len(values) < 3 || len(values) > 4 {
		fmt.Println("Invalid input, format is: ADD <KEY> <VALUE> [TABLE]")
	} else {
		name := default_database
		if len(values) == 4 {
			name = values[3]
			_, ok := database[name]
			if !ok {
				fmt.Println("Invalid table name, there is no table with name", name)
				return
			}
		}

		database[name][values[1]] = values[2]
		fmt.Println("Sucessfully added entry with key", values[1], "and value", values[2])
	}
}

func command_remove(values []string) {
	if len(values) < 2 || len(values) > 3 {
		fmt.Println("Invalid input, format is: REMOVE <KEY> [TABLE]")
	} else {
		name := default_database
		if len(values) == 3 {
			name = values[2]
			_, ok := database[name]
			if !ok {
				fmt.Println("Invalid table name, there is no table with name", name)
				return
			}
		}

		_, ok := database[name][values[1]]
		delete(database, values[1])
		if ok {
			fmt.Println("Sucessfully removed entry with key", values[1], "from database", name)
		} else {
			fmt.Println("Did not find entry with key", values[1], "in database", name)
		}
	}
}

func command_get(values []string) {
	if len(values) < 2 || len(values) > 3 {
		fmt.Println("Invalid input, format is: GET <KEY> [TABLE]")
	} else {
		name := default_database
		if len(values) == 3 {
			name = values[2]
			_, ok := database[name]
			if !ok {
				fmt.Println("Invalid table name, there is no table with name", name)
				return
			}
		}

		value, ok := database[name][values[1]]
		if ok {
			fmt.Println("Key", values[1], "was found in database", name, "with value", value)
		} else {
			fmt.Println("Key", values[1], "was not found in database", name)
		}
	}
}

func command_list(values []string) {
	if len(values) > 2 {
		fmt.Println("Invalid input, format is: LIST [TABLE]")
	} else {
		name := default_database
		if len(values) == 2 {
			name = values[1]
			_, ok := database[name]
			if !ok {
				fmt.Println("Invalid table name, there is no table with name", name)
				return
			}
		}

		fmt.Println("Database key and values:")
		for key, value := range database[name] {
			fmt.Println(key + ":" + value)
		}
	}
}

func command_add_table(values []string) {
	if len(values) != 2 {
		fmt.Println("Invalid input, format is: ADDTABLE <NAME>")
	} else {
		_, ok := database[values[1]]

		if ok {
			fmt.Println("Database with name", values[1], "already exists")
		} else {
			database[values[1]] = make(map[string]string)
			fmt.Println("Successfully added database with name", values[1])
		}
	}
}

func command_list_tables() {
	fmt.Println("List of tables:")
	for key := range database {
		fmt.Println(key)
	}
}

func load() {
	// Get all files
	directory, _ := os.Getwd()
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// Only .data files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".data") {
			continue
		}

		file_name := file.Name()[:len(file.Name())-len(".data")]

		// Check file
		file_location := directory + "\\" + file.Name()
		file, err := os.Open(file_location)

		if os.IsNotExist(err) {
			fmt.Println("No file found to load data from")
			return
		} else if err != nil {
			panic(err)
		}

		fmt.Println("Loading database vales from file", file_location)

		// Load data line for line
		scanner := bufio.NewScanner(file)
		database[file_name] = make(map[string]string)

		for scanner.Scan() {
			values := strings.Split(scanner.Text(), ":")
			if len(values) != 2 {
				fmt.Println("Line is not formatted correctly:", values)
				continue
			}

			database[file_name][values[0]] = values[1]
		}

		file.Close()
	}
}

func save() {
	// Create files for all tables
	for table_name, table_values := range database {
		// Create (override) file
		directory, _ := os.Getwd()
		file_location := directory + "\\" + table_name + ".data"
		file, err := os.Create(file_location)

		if err != nil {
			fmt.Println("Error when creating file", file_location)
		}

		// Write all database data
		writer := bufio.NewWriter(file)
		for key, value := range table_values {
			writer.WriteString(key + ":" + value + "\n")
		}

		writer.Flush()
		file.Close()
	}
}
