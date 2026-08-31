package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type server struct {
    clients map[net.Conn]any
    mutex   sync.Mutex
}

type client struct {
    connection net.Conn
    message		string
}

func main() {
	messages := make(chan client)
	listener, err := net.Listen("tcp", ":8080")
	srv := &server{clients: make(map[net.Conn]any)}
	
	if err != nil {
		fmt.Println("Something went wrong starting the server...")
		os.Exit(1)
	}

	fmt.Println("Server started ...")

	go broadcast(srv, messages)

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("Something went wrong connecting to the server...")
			fmt.Println("Please try again.")
		}

		go handleClient(connection, messages, srv)

		srv.mutex.Lock()
		srv.clients[connection] = true
		srv.mutex.Unlock()
	}
}
