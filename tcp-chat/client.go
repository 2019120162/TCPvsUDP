package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print(message)
		}
	}()

	input := bufio.NewReader(os.Stdin)
	for {
		text, _ := input.ReadString('\n')
		conn.Write([]byte(text))
	}
}
