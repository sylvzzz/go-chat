package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// handleClient speaks the chat protocol, whatever the transport is:
// it only sees a stream of strings + a way to write to itself.
func handleClient(c conn, msgs <-chan string, messages chan client, srv *server) {
	set := false // was the username the first message we got?
	username := ""

	for msg := range msgs {
		msg = strings.TrimSpace(msg)

		// first message is always the username
		if !set {
			if msg == "" {
				c.Close()
				return
			}
			username = msg
			set = true
			fmt.Printf("%v joined ...\n", username)
			messages <- client{connection: c, message: username + " joined the chat..."}
			continue
		}

		if msg == "exit" {
			break
		}

		message := username + ": " + msg
		fmt.Println(message)
		messages <- client{connection: c, message: message}
	}

	// connected but never gave a username, don't broadcast a fake leave
	if !set {
		c.Close()
		return
	}

	// broadcast the leave on every exit path (explicit exit or dropped connection)
	messages <- client{connection: c, message: username + " left the chat..."}
	fmt.Printf("%v left the chat ...\n", username)

	c.Close()

	// deleting connection from server
	srv.mutex.Lock()
	delete(srv.clients, c)
	srv.mutex.Unlock()
}

// tcpReader turns newline-separated TCP lines into the message stream
func tcpReader(connection net.Conn) <-chan string {
	msgs := make(chan string)
	go func() {
		scanner := bufio.NewScanner(connection)
		for scanner.Scan() {
			msgs <- scanner.Text()
		}
		close(msgs)
	}()
	return msgs
}
