package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleClient(connection net.Conn, messages chan string) {
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Client: %v\n", message)
		messages <- message
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	} else {
		fmt.Println("Client disconnected...")
	}

	connection.Close()
}