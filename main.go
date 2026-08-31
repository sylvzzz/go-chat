package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	clients := make(map[net.Conn]any)
	listener, err := net.Listen("tcp", ":8080")
	
	if err != nil {
		fmt.Println("Something went wrong starting the server...")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("Something went wrong connecting to the server...")
			fmt.Println("Please try again.")
		}

		go handleClient(connection)

		clients[connection] = true
	}
}
