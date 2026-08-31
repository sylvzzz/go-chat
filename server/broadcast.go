package main

import (
		"fmt"
		)

func broadcast(srv *server, messages chan client) {
	for {
		msg := <-messages
		srv.mutex.Lock()
		for connection := range srv.clients {
			if connection == msg.connection {
				continue
			}
			fmt.Fprintf(connection, "%v\n", msg)
		}
		srv.mutex.Unlock()
	}
}