package handler

import (
    "chat/server"
    "net"
)

func StartTCP(listener net.Listener) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go server.HandleTCPClient(conn)
    }
}