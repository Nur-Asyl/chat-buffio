package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

const (
	CONN_TEST_PORT_CLIENT = ":3334"
	CONN_TEST_TYPE_CLIENT = "tcp"
)

func Test_Client_MainConnection(t *testing.T) {
	t.Parallel()
	conn, err := net.Dial(CONN_TEST_TYPE_CLIENT, CONN_TEST_PORT_CLIENT)
	if err != nil {
		t.Skipf("Skipping test: unable to connect to %s: %s", CONN_PORT_CLIENT, err.Error())
		return
	}
	defer conn.Close()
}

func Test_Client_SendMessage(t *testing.T) {
	t.Parallel()
	conn, err := net.Dial(CONN_TYPE_CLIENT, CONN_PORT_CLIENT)
	if err != nil {
		t.Skipf("Skipping test: unable to connect to %s: %s", CONN_PORT_CLIENT, err.Error())
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString("Test message\n")
	writer.Flush()
}

func Test_Client_ReceiveMessage(t *testing.T) {
	t.Parallel()
	conn, err := net.Dial(CONN_TYPE_CLIENT, CONN_PORT_CLIENT)
	if err != nil {
		t.Skipf("Skipping test: unable to connect to %s: %s", CONN_PORT_CLIENT, err.Error())
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Error reading message: %s", err.Error())
		return
	}
	fmt.Println("Received message:", message)
}
