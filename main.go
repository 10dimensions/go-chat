package main

import (
    "chat/handler"
    "chat/server"
    "chat/ws"
    "embed"
    "fmt"
    "log"
    "net"
    "net/http"
)

//go:embed static
var staticFS embed.FS

func main() {
    // Start broadcaster goroutine
    go server.Broadcaster()

    // TCP server
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    fmt.Println("TCP Chat server started on :8080 (use nc localhost 8080 or Go client)")
    go handler.StartTCP(listener)

    // HTTP + WebSocket server on :8081
    http.HandleFunc("/ws", ws.HandleWS)

    // Serve the thin HTML5 client
    http.Handle("/", http.FileServer(http.FS(staticFS)))

    fmt.Println("Web interface available at http://localhost:8081")
    fmt.Println("WebSocket endpoint: ws://localhost:8081/ws")

    log.Fatal(http.ListenAndServe(":8081", nil))
}