package main

import (
	"net"
	"strings"
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
	this.server.BroadCast(this, " This user is online!")
}

func (this *User) Offline()  {
	//add user to OnlineMap (in server side, to keep record)
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//broadcast to user a success msg
	this.server.BroadCast(this, " This user is offline.")
}

func (this *User) SendMessage(msg string) {
	//this.C <-msg //send the message only to one user, not broadcaset. if use this, when force quit, can't print msg in client, why?
	this.conn.Write([]byte(msg+"\n")) //teaching material says so, but my way is better
}

func (this *User) ProcessMessage(msg string)  {
	if msg == "who" { //if user input "who", query current online user list
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "["+user.Addr+"]"+user.Name+": This user is currently online."
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg)>7 && msg[:7] == "rename|"{ //if user input "rename|Jack", rename this user to Jack
		newName := strings.Split(msg, "|")[1] //get "Jack"
		_, ok := this.server.OnlineMap[newName] //check if there already exists "Jack"
		if ok { //there already exists "Jack"
			this.SendMessage("The new name \" "+newName+" \" you input already exists, rename fails.")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("Your name has been updated to \" " + newName + "\".")
		}

	} else if len(msg)>3 && msg[:3] == "to|"{ //private message
		receiverName := strings.Split(msg, "|")[1] //get "Jack" the message receiver name
		if receiverName == "" || len(strings.Split(msg, "|"))<3 {
			this.SendMessage("Wrong format for private message.")
			return
		}

		receiverUser, ok := this.server.OnlineMap[receiverName]
		if ok { //if the receiver exists
			privateMessage := strings.Split(msg, "|")[2]
			if privateMessage == "" {
				this.SendMessage("You can't send empty message.")
				return
			}
			this.SendMessage("You sent a private message to " + receiverName + ": "+ privateMessage)
			receiverUser.SendMessage(this.Name + " sent a private message to you: " + privateMessage)
		} else { //if the receiver doesn't exist
			this.SendMessage("The receiver of your private message doesn't exist.")
			return
		}

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