package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Room_CreateRoom(t *testing.T) {
	t.Parallel()

	cs := NewChatServer()
	roomName := "testRoom"
	room := cs.CreateRoom(roomName)

	assert.NotNil(t, room)
	assert.Equal(t, roomName, room.name)
	assert.NotNil(t, room.clients)
	assert.NotNil(t, room.logFile)
	assert.NotNil(t, room.logWriter)

	// Clean up
	delete(cs.Rooms, roomName)
	os.Remove(room.name + "_log.txt")
}
