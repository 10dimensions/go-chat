package server

import (
    "bufio"
    "fmt"
    "strings"
)

func HandleTCPClient(conn net.Conn) {
    defer conn.Close()

    conn.Write([]byte("Enter your username: "))
    reader := bufio.NewReader(conn)
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    if name == "" {
        name = "Anonymous"
    }

    client := Client{Conn: conn, Name: name}

    mutex.Lock()
    TCPclients[conn] = client
    mutex.Unlock()

    broadcast <- fmt.Sprintf("%s has joined the chat!", name)

    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            mutex.Lock()
            delete(TCPclients, conn)
            mutex.Unlock()
            broadcast <- fmt.Sprintf("%s has left the chat.", name)
            return
        }

        msg = strings.TrimSpace(msg)
        if msg != "" {
            broadcast <- fmt.Sprintf("%s: %s", name, msg)
        }
    }
}