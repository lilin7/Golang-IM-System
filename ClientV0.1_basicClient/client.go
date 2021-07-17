package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// instantiate a client object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil { //if fail connecting server
		fmt.Println("net.Dial error:", err)
		return nil
	}
	// if succeed connecting server
	client.conn = conn
	return client
}

func main() {
	//hard code ip and port for now
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("Failed conneting server!")
		return
	}

	fmt.Println("Succeed connecting server!")

	//block main routine
	select {

	}
}
