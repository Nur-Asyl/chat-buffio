package main

import (
	"bufio"
	"ex1/internal"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":3334")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 3334")

	chatServer := internal.NewChatServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}

		client := &internal.Client{
			Conn:   conn,
			Reader: bufio.NewReader(conn),
			Writer: bufio.NewWriter(conn),
		}

		go chatServer.CreateRoom("default").HandleClient(client, chatServer)
	}
}
