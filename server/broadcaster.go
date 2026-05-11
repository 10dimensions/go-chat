package server

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	TCPclients = make(map[net.Conn]Client)
	WSclients  = make(map[*websocket.Conn]WSClient)
	Broadcast  = make(chan string) // Exported
	mutex      = sync.Mutex{}      // Keep unexported (we use functions if needed)
)

func Broadcaster() {
	for msg := range Broadcast {
		mutex.Lock()

		// TCP clients
		for conn := range TCPclients {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				conn.Close()
				delete(TCPclients, conn)
			}
		}

		// WebSocket clients
		for conn := range WSclients {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				conn.Close()
				delete(WSclients, conn)
			}
		}

		mutex.Unlock()
	}
}
