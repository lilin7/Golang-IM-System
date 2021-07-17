package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

	fmt.Scanln(&flag) //take user keyboard input, and write to the address of variable flag

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
			client.SendPublicMessage()
			break
		case 2:
			fmt.Println("You have selected: Send private message")
			client.SendPrivateMessage()
			break
		case 3:
			fmt.Println("You have selected: Rename yourself")
			client.ChangeUserName()
			break
		}
	}
}

func (client *Client) ChangeUserName() bool {

	fmt.Println("Please input a new name for yourself:")
	fmt.Scanln(&client.Name) //take user keyboard input, and write to the address of variable client.Name

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

// this will be a goroutine to process server response and print to client console
func (client *Client) ProcessServerResponse() {
	io.Copy(os.Stdout, client.conn) // once client.conn passes data, copy to stdout, forever block and wait for client.conn

/*	//same as:
	for {
		buf := make([]byte, 4096)
		client.conn.Read(buf)
		fmt.Println(buf)
	}*/
}

func (client *Client) SendPublicMessage() {
	var chatMsg string

	fmt.Println("Please input public message to be sent: (input \"exit\" to quit)")
	fmt.Scanln(&chatMsg) //take user keyboard input, and write to the address of variable chatMsg

	for chatMsg != "exit" {
		if len(chatMsg) != 0 { //if input is not empty
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("Please input public message to be sent: (input \"exit\" to quit)")
		fmt.Scanln(&chatMsg)
	}

}

// query the list of online users
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

//send private message to another client
func (client *Client) SendPrivateMessage() {
	var receiverName string
	var privateMessage string

	client.SelectUsers()
	fmt.Println("Please input a name to send a private message to: (input \"exit\" to quit)")
	fmt.Scanln(&receiverName) //take user keyboard input, and write to the address of variable receiverName

	for receiverName != "exit" {
		fmt.Println("Please input the private message: (input \"exit\" to quit)")
		fmt.Scanln(&privateMessage) //take user keyboard input, and write to the address of variable privateMessage

		for privateMessage != "exit" {
			if len(privateMessage) != 0 {
				sendMsg := "to|" + receiverName + "|" + privateMessage + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			privateMessage = ""
			fmt.Println("Please input the private message: (input \"exit\" to quit)")
			fmt.Scanln(&privateMessage)
		}

		client.SelectUsers()
		fmt.Println("Please input a name to send a private message to: (input \"exit\" to quit)")
		fmt.Scanln(&receiverName)
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

	//a goroutine to process server response and print to client console, forever block and wait for client.conn
	go client.ProcessServerResponse()

	client.Run() //major business after connecting to server
}
