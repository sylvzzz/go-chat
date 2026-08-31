package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	)

func handleClient(connection net.Conn, messages chan string, srv *server) {
	scanner := bufio.NewScanner(connection)

	fmt.Fprint(connection, "Enter your username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Printf("%v joined ...", username)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			messages <- username + " left the chat...\n"
			break
		}

		message := username + ": " + input
		fmt.Println( message)
		messages <- message
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading: ", err)
	} else {
		fmt.Printf("%v left the chat ...", username)
	}

	connection.Close()

	// deleting connection from server
	srv.mutex.Lock()
	delete(srv.clients, connection)
	srv.mutex.Unlock()
}