package internal

import (
	"bufio"
	"net"
)

type Client struct {
	Conn   net.Conn
	Name   string
	Room   *Room
	Reader *bufio.Reader
	Writer *bufio.Writer
}

type ChatServer struct {
	Rooms map[string]*Room
}

func NewChatServer() *ChatServer {
	return &ChatServer{Rooms: make(map[string]*Room)}
}

func (cs *ChatServer) CreateRoom(name string) *Room {
	room := &Room{
		name:    name,
		clients: make(map[net.Conn]*Client),
	}
	room.createLogFile()
	cs.Rooms[name] = room
	return room
}

func (cs *ChatServer) GetRoom(name string) *Room {
	return cs.Rooms[name]
}
