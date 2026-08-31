package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	)

func handleClient(connection net.Conn, messages chan string) {
	scanner := bufio.NewScanner(connection)

	fmt.Fprint(connection, "Enter your username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	for scanner.Scan() {
		message := username + ": " + scanner.Text()
		fmt.Println( message)
		messages <- message
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	} else {
		fmt.Println("Client disconnected...")
	}

	connection.Close()
}