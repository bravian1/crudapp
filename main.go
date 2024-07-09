package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func scanText() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimSpace(text)
}

func main() {
	db := DBConn()
	defer db.conn.Close()
	db.maybeCreateTable()
	println("# Welcome to student app")
	println()

	for {
		fmt.Println("\nEnter your choice:")
		fmt.Println("1. Add user")
		fmt.Println("2. Show all users")
		fmt.Println("3. Get user by ID")
		fmt.Println("4. Edit user")
		fmt.Println("5. Delete user")
		fmt.Println("6. Exit")
		choice := scanText()

		switch choice {
		case "1":
			fmt.Println("Enter username: ")
			name := scanText()
			fmt.Println("Enter age: ")
			age := scanText()
			err := db.addUser(name, age)
			if err != nil {
				fmt.Printf("Error adding user: %v\n", err)
			} else {
				fmt.Println("User added successfully")
			}
		case "2":
			students, err := db.getAll()
			if err != nil {
				fmt.Printf("Error getting users: %v\n", err)
			} else {
				for _, student := range students {
					fmt.Println(student.String())
				}
			}
		case "3":
			fmt.Println("Enter user ID: ")
			id, err := strconv.Atoi(scanText())
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}
			student, err := db.GetById(id)
			if err != nil {
				fmt.Printf("Error getting user: %v\n", err)
			} else {
				fmt.Println(student.String())
			}
		case "4":
			fmt.Println("Enter user ID: ")
			id, err := strconv.Atoi(scanText())
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}
			fmt.Println("Enter updated name: ")
			name := scanText()
			fmt.Println("Enter updated age: ")
			age := scanText()
			err = db.editUser(id, name, age)
			if err != nil {
				fmt.Printf("Error editing user: %v\n", err)
			} else {
				fmt.Println("User edited successfully")
			}
		case "5":
			fmt.Println("Enter user ID: ")
			id, err := strconv.Atoi(scanText())
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}
			err = db.deleteUser(id)
			if err != nil {
				fmt.Printf("Error deleting user: %v\n", err)
			} else {
				fmt.Println("User deleted successfully")
			}
		case "6":
			fmt.Println("Exiting the application. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}