package server

import "net"

type Client struct {
    Conn net.Conn
    Name string
}

type WSClient struct {
    Conn *websocket.Conn
    Name string
}