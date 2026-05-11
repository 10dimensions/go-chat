package ws

import (
	"chat/server"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Simple username handling (sent as first message)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}
	name := string(msg)
	if name == "" {
		name = "Anonymous"
	}

	client := server.WSClient{Conn: conn, Name: name}

	server.mutex.Lock()
	server.WSclients[conn] = client
	server.mutex.Unlock()

	server.Broadcast <- fmt.Sprintf("%s has joined the chat!", name)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			server.mutex.Lock()
			delete(server.WSclients, conn)
			server.mutex.Unlock()
			server.Broadcast <- fmt.Sprintf("%s has left the chat.", name)
			return
		}

		text := string(message)
		if text != "" {
			server.Broadcast <- fmt.Sprintf("%s: %s", name, text)
		}
	}
}
