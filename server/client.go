package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	)


func handleClient(connection net.Conn, messages chan client, srv *server) {
	scanner := bufio.NewScanner(connection)

	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if username == "" {
		connection.Close()
		return
	}

	fmt.Printf("%v joined ...\n", username)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			messages <- client{connection: connection, message: username + " left the chat...\n"}
			break
		}

		message := username + ": " + input
		fmt.Println(message)
		messages <- client{connection: connection, message: message}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading: ", err)
	} else {
		fmt.Printf("%v left the chat ...\n", username)
	}

	connection.Close()

	// deleting connection from server
	srv.mutex.Lock()
	delete(srv.clients, connection)
	srv.mutex.Unlock()
}