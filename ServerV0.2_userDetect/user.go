package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string //channel to be used to send msg to client
	conn net.Conn //the connection communicate with server
}

//create a user interface, returns a server pointer
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String() //get the user address from connection

	//create an instance of Object "User", assign the values get from conn, then return (a pointer to this User object)
	user := &User{
		Name : userAddr,
		Addr : userAddr,
		C: make(chan string), //make a channel C, to be used to send msg to client
		conn: conn,
	}
	//start a goroutine to listen the channel of current user, once there is msg, send to client
	go user.ListenMessage()

	return user
}

//keep listening to current User channel, once there is message, send to client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte (msg + "\n"))
	}

}