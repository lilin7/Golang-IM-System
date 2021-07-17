//server
package main

import (
	"fmt"
	"io"
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
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
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
	user := NewUser(conn, this)

	//add user to OnlineMap (in server side, to keep record), broadcast to user a success msg
	user.Online()

	fmt.Println("[" + user.Addr + "]" + user.Name + ": " + "This user is online!")


	//receive msg sent by users
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) //get user input, n is the length of msg
			if n==0 { //if there is no msg sent by a user, means this user is offline
				user.Offline() //delete from usermap, broadcast
				fmt.Println("[" + user.Addr + "]" + user.Name + "This user is offline.")
				return
			}

			if err != nil && err != io.EOF { //if there is error, or not end of file
				fmt.Println("Conn Read error:", err)
				return
			}

			//get user input. because end of buf is \n, remove it, and convert from byte[] to string
			msg := string(buf[:n-1])
			user.ProcessMessage(msg) //broadcast
			fmt.Println("[" + user.Addr + "]" + user.Name + " sends message: " + msg)

		}
	}()

	//let thread to wait?????
	select {

	}
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
