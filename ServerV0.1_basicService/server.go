//server
package main

import (
	"fmt"
	"net"
)

//a Server class, it has 2 attributes
type Server struct {
	Ip string
	Port int
}

//create a server interface, returns a server pointer
func NewServer(ip string, port int) *Server {
	//create an instance of Object "Server", assign the address of it to varibale "server", then return (a pointer to this Server object)
	server := &Server{
		Ip : ip,
		Port : port,
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	//
	fmt.Println("Connection built successfully!")

}

//start server (bind this method to Server object)
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen error:", err)
	}

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
