package main

import (
		"fmt"
		)

func broadcast(srv *server, messages chan string) {
	for {
		msg := <-messages
		srv.mutex.Lock()
		for connection := range srv.clients {
			fmt.Fprintf(connection, "%v\n", msg)
		}
		srv.mutex.Unlock()
	}
}