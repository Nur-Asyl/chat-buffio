package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChatServer(t *testing.T) {
	t.Parallel()
	// Arrange
	cs := NewChatServer()

	// Assert
	assert.NotNil(t, cs, "Expected a non-nil ChatServer instance, got nil")
	assert.Empty(t, cs.Rooms, "Expected no rooms to be created initially")
}

func TestCreateRoom(t *testing.T) {
	t.Parallel()
	// Arrange
	cs := NewChatServer()
	roomName := "testRoom"

	// Act
	room := cs.CreateRoom(roomName)

	// Assert
	assert.NotNil(t, room, "Expected a non-nil Room instance, got nil")
	assert.Contains(t, cs.Rooms, roomName, "Expected room '%s' to be created, but it's not in the map", roomName)
}

func TestGetRoom(t *testing.T) {
	t.Parallel()
	// Arrange
	cs := NewChatServer()
	roomName := "testRoom"
	room := cs.CreateRoom(roomName)

	// Act
	retrievedRoom := cs.GetRoom(roomName)

	// Assert
	assert.Equal(t, room, retrievedRoom, "Expected to get room '%s', got different room", roomName)
}
