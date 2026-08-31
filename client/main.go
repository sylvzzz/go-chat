package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	connection, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting to the server please try again...")
		os.Exit(1)
	}
	
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Fprintf(connection, "%v\n", username)

	if err := scanner.Err(); err != nil {
		fmt.Println("Something went wrong please try again ...")
	}

	fmt.Printf("Me: ")

	go func() {
		serverScanner := bufio.NewScanner(connection)

		if err := serverScanner.Err(); err != nil {
			fmt.Println("Something went wrong please try again ...")
		}

		for serverScanner.Scan() {
			fmt.Println(serverScanner.Text())
		}
	}()

	for scanner.Scan() {
		fmt.Printf("Me: ")
		message := scanner.Text()
		if message == "exit" {
			break
		}
		fmt.Fprintf(connection, "%v\n", message)
	}
}