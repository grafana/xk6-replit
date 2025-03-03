package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// Response struct for JSON parsing
type Response struct {
	Status  string `json:"status"`
	Output  string `json:"output,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

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
		input = strings.TrimSpace(input)

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected. Exiting REPL.")
			return
		}

		var res Response
		err = json.Unmarshal([]byte(response), &res)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			continue
		}

		// Handle response based on status
		if res.Status == "ok" {
			fmt.Println(res.Output)
		} else if res.Status == "error" {
			fmt.Println("Error:", res.Error) // Now correctly displaying the Error field
		} else if res.Status == "exit" {
			fmt.Println(res.Message)
			return
		}
	}
}
