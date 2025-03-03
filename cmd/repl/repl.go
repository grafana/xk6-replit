package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
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
		color.Red("Failed to connect to REPL server: %v", err)
		return
	}
	defer conn.Close()

	color.Cyan("Connected to k6 REPL. Type 'exit' to quit.")

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			color.Red("Error sending command: %v", err)
			return
		}

		response, err := serverReader.ReadString('\n')
		if err != nil {
			color.Red("Server disconnected. Exiting REPL.")
			return
		}

		var res Response
		err = json.Unmarshal([]byte(response), &res)
		if err != nil {
			color.Red("Error parsing response: %v | Raw: %s", err, response)
			continue
		}

		// Handle response based on status with color output
		switch res.Status {
		case "ok":
			color.Green(res.Output) // Green for success
		case "error":
			color.Red("Error: %s", res.Error) // Red for errors
		case "exit":
			color.Yellow(res.Message) // Yellow for exit message
			return
		}
	}
}
