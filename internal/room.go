package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Room struct {
	name      string
	clients   map[net.Conn]*Client
	logFile   *os.File
	logWriter *bufio.Writer
}

func (cs *ChatServer) createRoom(name string) *Room {
	room := &Room{
		name:    name,
		clients: make(map[net.Conn]*Client),
	}
	room.createLogFile()
	cs.Rooms[name] = room
	return room
}

func (cs *ChatServer) getRoom(name string) *Room {
	return cs.Rooms[name]
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
	room.logWriter.WriteString(fmt.Sprintf("[%s] %s: %s\n", now, client.Name, message))
	room.logWriter.Flush()
}

func (room *Room) broadcast(sender *Client, message string) {
	for conn, client := range room.clients {
		if conn != sender.Conn {
			client.Writer.WriteString(sender.Name + ": " + message + "\n")
			client.Writer.Flush()
		}
	}
	room.logMessage(sender, message)
}

func (room *Room) join(client *Client) {
	client.Room = room
	room.clients[client.Conn] = client
	client.Writer.WriteString("Joined room: " + room.name + "\n")
	client.Writer.Flush()
}

func (room *Room) leave(client *Client) {
	delete(room.clients, client.Conn)
	client.Room = nil
	client.Writer.WriteString("Left room: " + room.name + "\n")
	client.Writer.Flush()
}

func (room *Room) HandleClient(client *Client, cs *ChatServer) {
	defer client.Conn.Close()
	fmt.Printf("New client connected: %s\n", client.Conn.RemoteAddr().String())

	client.Writer.WriteString("Welcome to the chat room!\nEnter your name: ")
	client.Writer.Flush()

	client.Name, _ = client.Reader.ReadString('\n')
	client.Name = strings.TrimSpace(client.Name)

	client.Writer.WriteString("Welcome!!! " + client.Name + "!\n")
	client.Writer.WriteString("Use the /help command to see available commands!!! \\'0'/ \n")
	client.Writer.Flush()

	for {
		message, err := client.Reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "/create") {
			roomName := strings.TrimSpace(strings.TrimPrefix(message, "/create"))
			cs.createRoom(roomName)
			client.Writer.WriteString("Room created: " + roomName + "\n")
			client.Writer.Flush()
		} else if strings.HasPrefix(message, "/join") {
			roomName := strings.TrimSpace(strings.TrimPrefix(message, "/join"))
			targetRoom := cs.getRoom(roomName)
			if targetRoom == nil {
				client.Writer.WriteString("Room does not exist: " + roomName + "\n")
				client.Writer.Flush()
			} else {
				if client.Room != nil {
					client.Room.leave(client)
				}
				targetRoom.join(client)
				room = targetRoom
			}
		} else if message == "/exit" {
			break
		} else if message == "/help" {
			client.Writer.WriteString("Available commands:\n")
			client.Writer.WriteString("/create [room_name] - Create a new room.\n")
			client.Writer.WriteString("/join [room_name] - Join an existing room.\n")
			client.Writer.WriteString("/exit - Exit the chat.\n")
			client.Writer.Flush()
		} else {
			if room != nil {
				now := time.Now().Format("2006-01-02 15:04:05")
				room.broadcast(client, "["+now+"] "+message)
			} else {
				client.Writer.WriteString("You are not in room. Use /join to enter room.\n")
				client.Writer.Flush()
			}
		}
	}

	if room != nil {
		room.leave(client)
	}
	fmt.Printf("%s has disconnected\n", client.Name)
}
