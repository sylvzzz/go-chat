package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type server struct {
	clients map[conn]bool
	mutex   sync.Mutex
}

type conn interface {
	WriteMessage(msg string)
	Close() error
}

type tcpConn struct {
	net.Conn
}

func (c *tcpConn) WriteMessage(msg string) {
	fmt.Fprintf(c.Conn, "%v\n", msg)
}

type client struct {
	connection conn // was net.Conn, is the interface now
	message    string
}

func main() {
	messages := make(chan client)
	listener, err := net.Listen("tcp", ":8080")
	srv := &server{clients: make(map[conn]bool)}

	if err != nil {
		fmt.Println("Something went wrong starting the server...")
		os.Exit(1)
	}

	fmt.Println("Server started on port 8080...")

	go broadcast(srv, messages)
	go startWS(srv, messages)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Something went wrong connecting to the server...")
			fmt.Println("Please try again.")
			continue
		}

		// raw net.Conn gets a WriteMessage + a reader, so it fits the hub
		c := &tcpConn{Conn: connection}
		go handleClient(c, tcpReader(connection), messages, srv)

		srv.mutex.Lock()
		srv.clients[c] = true
		srv.mutex.Unlock()
	}
}
