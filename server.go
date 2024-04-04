package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn   net.Conn
	name   string
	room   *Room
	reader *bufio.Reader
	writer *bufio.Writer
}

type Room struct {
	name    string
	clients map[net.Conn]*Client
}

type ChatServer struct {
	rooms map[string]*Room
}

func NewChatServer() *ChatServer {
	return &ChatServer{rooms: make(map[string]*Room)}
}

func (cs *ChatServer) createRoom(name string) *Room {
	room := &Room{name: name, clients: make(map[net.Conn]*Client)}
	cs.rooms[name] = room
	return room
}

func (cs *ChatServer) getRoom(name string) *Room {
	return cs.rooms[name]
}

func (room *Room) broadcast(sender *Client, message string) {
	for conn, client := range room.clients {
		if conn != sender.conn {
			client.writer.WriteString(sender.name + ": " + message + "\n")
			client.writer.Flush()
		}
	}
}
func (room *Room) join(client *Client) {
	client.room = room
	room.clients[client.conn] = client
	client.writer.WriteString("Joined room: " + room.name + "\n")
	client.writer.Flush()

	for _, c := range room.clients {
		if c != client {
			c.writer.WriteString(client.name + " has joined the room.\n")
			c.writer.Flush()
		}
	}
}

func (room *Room) leave(client *Client) {
	for _, c := range room.clients {
		if c != client {
			c.writer.WriteString(client.name + " has left the room.\n")
			c.writer.Flush()
		}
	}
	delete(room.clients, client.conn)
	client.room = nil
	client.writer.WriteString("Left room: " + room.name + "\n")
	client.writer.Flush()
}
func (room *Room) handleClient(client *Client, cs *ChatServer) {
	defer client.conn.Close()
	fmt.Printf("New client connected: %s\n", client.conn.RemoteAddr().String())

	client.writer.WriteString("Welcome to the chat room!\nEnter your name: ")
	client.writer.Flush()

	client.name, _ = client.reader.ReadString('\n')
	client.name = strings.TrimSpace(client.name)

	client.writer.WriteString("Welcome, " + client.name + "!\n")
	client.writer.WriteString("Use the /help command to see available commands!!!!! \\'0'/\n")
	client.writer.Flush()

	for {
		message, err := client.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "/create") {
			roomName := strings.TrimSpace(strings.TrimPrefix(message, "/create"))
			cs.createRoom(roomName)
			client.writer.WriteString("Room created: " + roomName + "\n")
			client.writer.Flush()
		} else if strings.HasPrefix(message, "/join") {
			roomName := strings.TrimSpace(strings.TrimPrefix(message, "/join"))
			targetRoom := cs.getRoom(roomName)
			if targetRoom == nil {
				client.writer.WriteString("Room does not exist: " + roomName + "\n")
				client.writer.Flush()
			} else {
				if client.room != nil {
					client.room.leave(client)
				}
				targetRoom.join(client)
				room = targetRoom
			}
		} else if message == "/exit" {
			break
		} else if message == "/help" {
			client.writer.WriteString("Available commands:\n")
			client.writer.WriteString("/create [room_name] - Create new room.\n")
			client.writer.WriteString("/join [room_name] - Join an existing room.\n")
			client.writer.WriteString("/exit - Exit the chat.\n")
			client.writer.Flush()
		} else {
			if room != nil {
				now := time.Now().Format("2006-01-02 15:04:05")
				room.broadcast(client, "["+now+"] "+message)
			} else {
				client.writer.WriteString("You are not in a room. Use /join to enter room.\n")
				client.writer.Flush()
			}
		}
	}

	if room != nil {
		room.leave(client)
	}
	fmt.Printf("%s has disconnected\n", client.name)
}

const (
	CONN_PORT = ":3334"
	CONN_TYPE = "tcp"
)

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on portt 3334")

	chatServer := NewChatServer()

	defaultRoom := chatServer.createRoom("default")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}

		client := &Client{
			conn:   conn,
			reader: bufio.NewReader(conn),
			writer: bufio.NewWriter(conn),
		}

		go defaultRoom.handleClient(client, chatServer)
	}
}
