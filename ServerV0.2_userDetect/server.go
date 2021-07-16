//server
package main

import (
	"fmt"
	"net"
	"sync"
)

//a Server class, it has 2 attributes
type Server struct {
	Ip string
	Port int

	// a map to store all online users
	OnlineMap map[string] *User
	mapLock sync.RWMutex //make it multi-thread safe

	//a channel to send out msg
	Message chan string
}

//create a server interface, returns a server pointer
func NewServer(ip string, port int) *Server {
	//create an instance of Object "Server", assign the address of it to varibale "server", then return (a pointer to this Server object)
	server := &Server{
		Ip : ip,
		Port : port,
		OnlineMap: make(map[string] *User), //keep a map to store user info in the server
		Message: make(chan string), //keep a channel to send out msg
	}

	return server
}

// a method of Server, to broadcast a message to all online users
func (this *Server) BroadCast (user *User, msg string) { //sent by which user, content of msg
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg //broadcast this message to all online users
}

// once there is a broadcast message, loop over all online users and broadcast it out
func (this *Server) ListenMessage()  {
	for {
		msg := <- this.Message

		this.mapLock.Lock()

		for _, cli := range this.OnlineMap {
			cli.C <- msg //send msg to the channel C of all users
		}

		this.mapLock.Unlock()

	}
	
}

func (this *Server) Handler(conn net.Conn) {

	//once there is a incoming connection (conn), means one user is online, create this User object in server
	user := NewUser(conn)

	//add user to OnlineMap (in server side, to keep record)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//broadcast to user a success msg
	this.BroadCast(user, "User is online!")
}

//start server (bind this method to Server object)
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen error:", err)
	}

	go this.ListenMessage() //inside is a for loop to listen to "Message" (from broadcast) forever

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener error:", err)
			continue
		}

		//do handler
		//create a goroutine to handle the received conn, and Server goes into next for loop to accept next incoming conn
		go this.Handler(conn)
	}


	//close listen socket
	//will execute after the execution of the whole method
	defer listener.Close()

}
