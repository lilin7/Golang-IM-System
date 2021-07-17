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

	flag int //to record client input (choice on menu)
}

func NewClient(serverIp string, serverPort int) *Client {
	// instantiate a client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: -1, //as long as not 0, because 0 means exit from menu
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. Send public message")
	fmt.Println("2. Send private message")
	fmt.Println("3. Rename yourself")
	fmt.Println("0. Quit")

	fmt.Scanln(&flag) //take user keyboard input

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Wrong number, please enter again!")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 { //if client doesn't want to exit, go in this for loop
		for client.menu() != true { // client input wrong number (nothing from menu)
		}
		//client input 1 or 2 or 3:
		switch client.flag {
		case 1:
			fmt.Println("You have selected: Send public message")
			break
		case 2:
			fmt.Println("You have selected: Send private message")
			break
		case 3:
			fmt.Println("You have selected: Rename yourself")
			break
		}
	}
}

func main() {
	flag.Parse() //enable parsing command line like "./client -ip 127.0.0.1 -port 8888"
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Failed conneting server!")
		return
	}

	fmt.Println("Succeed connecting server!")

	client.Run() //major business after connecting to server
}
