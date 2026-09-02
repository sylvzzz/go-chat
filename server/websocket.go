package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// wsConn is the websocket twin of tcpConn, so the hub treats both the same
type wsConn struct {
	w *websocket.Conn
}

func (c *wsConn) WriteMessage(msg string) {
	c.w.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *wsConn) Close() error {
	return c.w.Close()
}

// wsReader feeds websocket frames into the same message stream as tcpReader
func wsReader(w *websocket.Conn) <-chan string {
	msgs := make(chan string)
	go func() {
		for {
			_, msg, err := w.ReadMessage()
			if err != nil {
				close(msgs)
				return
			}
			msgs <- string(msg)
		}
	}()
	return msgs
}

func startWS(srv *server, messages chan client) {
	upgrader := websocket.Upgrader{
		// dev-only: the Vite proxy also avoids CORS, drop for prod
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// same shape as the TCP accept loop, so handleClient is shared
		wc := &wsConn{w: c}
		go handleClient(wc, wsReader(c), messages, srv)

		srv.mutex.Lock()
		srv.clients[wc] = true
		srv.mutex.Unlock()
	})

	// tcp chomps :8080, so HTTP/ws lives next door
	fmt.Println("WebSocket server started on port 8081 ...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Something went wrong starting the ws server:", err)
	}
}
