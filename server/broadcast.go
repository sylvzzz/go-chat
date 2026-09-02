package main

func broadcast(srv *server, messages chan client) {
	for {
		msg := <-messages
		srv.mutex.Lock()
		for connection := range srv.clients {
			if connection == msg.connection {
				continue
			}
			// transport-agnostic send, TCP and WS both implement it
			connection.WriteMessage(msg.message)
		}
		srv.mutex.Unlock()
	}
}
