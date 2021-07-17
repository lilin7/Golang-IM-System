package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C chan string //channel to be used to send msg to client
	conn net.Conn //the connection communicate with server

	server *Server //this user belongs to which server
}

//create a user interface, returns a server pointer
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String() //get the user address from connection

	//create an instance of Object "User", assign the values get from conn, then return (a pointer to this User object)
	user := &User{
		Name : userAddr,
		Addr : userAddr,
		C: make(chan string), //make a channel C, to be used to send msg to client
		conn: conn,

		server: server,
	}
	//start a goroutine to listen the channel of current user, once there is msg, send to client
	go user.ListenMessage()

	return user
}

func (this *User) Online()  {
	//add user to OnlineMap (in server side, to keep record)
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//broadcast to user a success msg
	this.server.BroadCast(this, "This user is online!")
}

func (this *User) Offline()  {
	//add user to OnlineMap (in server side, to keep record)
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//broadcast to user a success msg
	this.server.BroadCast(this, "This user is offline.")
}

func (this *User) SendMessage(msg string) {
	this.C <-msg //send the message only to one user, not broadcaset. my way
	//this.conn.Write([]byte(msg)) //teaching material says so, but my way is better
}

func (this *User) ProcessMessage(msg string)  {
	if msg == "who" { //if user input "who", query current online user list
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "["+user.Addr+"]"+user.Name+": This user is currently online."
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else {
		this.server.BroadCast(this, msg)
	}

}

//keep listening to current User channel, once there is message, send to client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte (msg + "\n"))
	}

}