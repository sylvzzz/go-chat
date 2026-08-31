package main

import (
	"fmt"
	"net"
	"bufio"
)

func handleClient(connection net.Conn) {
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Client: %v\n", message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	} else {
		fmt.Println("Client disconnected...")
	}

	connection.Close()
}