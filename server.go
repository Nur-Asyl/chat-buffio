package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	name      string
	clients   map[net.Conn]*Client
	logFile   *os.File
	logWriter *bufio.Writer
}

type ChatServer struct {
	rooms map[string]*Room
}

func NewChatServer() *ChatServer {
	return &ChatServer{rooms: make(map[string]*Room)}
}

func (cs *ChatServer) createRoom(name string) *Room {
	room := &Room{
		name:    name,
		clients: make(map[net.Conn]*Client),
	}
	room.createLogFile()
	cs.rooms[name] = room
	return room
}

func (cs *ChatServer) getRoom(name string) *Room {
	return cs.rooms[name]
}

func (room *Room) createLogFile() {
	fileName := room.name + "_log.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	room.logFile = file
	room.logWriter = bufio.NewWriter(file)
}

func (room *Room) logMessage(client *Client, message string) {
	if room.logWriter == nil {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	room.logWriter.WriteString(fmt.Sprintf("[%s] %s: %s\n", now, client.name, message))
	room.logWriter.Flush()
}

func (room *Room) broadcast(sender *Client, message string) {
	for conn, client := range room.clients {
		if conn != sender.conn {
			client.writer.WriteString(sender.name + ": " + message + "\n")
			client.writer.Flush()
		}
	}
	room.logMessage(sender, message)
}

func (room *Room) join(client *Client) {
	client.room = room
	room.clients[client.conn] = client
	client.writer.WriteString("Joined room: " + room.name + "\n")
	client.writer.Flush()
}

func (room *Room) leave(client *Client) {
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

	client.writer.WriteString("Welcome!!! " + client.name + "!\n")
	client.writer.WriteString("Use the /help command to see available commands!!! \\'0'/ \n")
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
			client.writer.WriteString("/create [room_name] - Create a new room.\n")
			client.writer.WriteString("/join [room_name] - Join an existing room.\n")
			client.writer.WriteString("/exit - Exit the chat.\n")
			client.writer.Flush()
		} else {
			if room != nil {
				now := time.Now().Format("2006-01-02 15:04:05")
				room.broadcast(client, "["+now+"] "+message)
			} else {
				client.writer.WriteString("You are not in room. Use /join to enter room.\n")
				client.writer.Flush()
			}
		}
	}

	if room != nil {
		room.leave(client)
	}
	fmt.Printf("%s has disconnected\n", client.name)
}

func main() {
	listener, err := net.Listen("tcp", ":3334")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 3334")

	chatServer := NewChatServer()

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

		go chatServer.createRoom("default").handleClient(client, chatServer)
	}
}
