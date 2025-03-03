package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8089")
	if err != nil {
		fmt.Println("Failed to connect to REPL server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to k6 REPL. Type 'exit' to quit.")

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		response, _ := serverReader.ReadString('\n')
		fmt.Print(response)

		if input == "exit\n" {
			return
		}
	}
}
