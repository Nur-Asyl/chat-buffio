package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3334")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	writer.WriteString(name + "\n")
	writer.Flush()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	for {
		message, _ := reader.ReadString('\n')
		writer.WriteString(message)
		writer.Flush()
	}
}
