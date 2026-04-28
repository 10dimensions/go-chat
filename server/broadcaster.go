package server

import (
    "fmt"
    "sync"
)

var (
    TCPclients   = make(map[net.Conn]Client)
    WSclients    = make(map[*websocket.Conn]WSClient)
    broadcast    = make(chan string)
    mutex        = sync.Mutex{}
)

func Broadcaster() {
    for msg := range broadcast {
        mutex.Lock()
        // Broadcast to TCP clients
        for conn := range TCPclients {
            _, err := conn.Write([]byte(msg + "\n"))
            if err != nil {
                conn.Close()
                delete(TCPclients, conn)
            }
        }
        // Broadcast to WebSocket clients
        for conn, client := range WSclients {
            err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
            if err != nil {
                conn.Close()
                delete(WSclients, conn)
            }
        }
        mutex.Unlock()
    }
}