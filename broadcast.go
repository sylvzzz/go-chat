package main

import (
		"fmt"
		"net"
		)

func broadcast(clients map[net.Conn]any, messages chan string) {
	for {
		msg := <-messages
		for connection := range clients {
			fmt.Fprintf(connection, "%v: %v\n", msg)
		}
	}
}