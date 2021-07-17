package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// instantiate a client object
	client := &Client{
		ServerIp: serverIp,
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

var serverIp string
var serverPort int
//parse command line e.g.: "./client -ip 127.0.0.1 -port 8888"
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Configure server IP address (default is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Configure server port (default is 8888)")
}

func main() {
	flag.Parse() //enable parsing command line like "./client -ip 127.0.0.1 -port 8888"
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Failed conneting server!")
		return
	}

	fmt.Println("Succeed connecting server!")

	//block main routine
	select {

	}
}
