package replit

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

// API is the exposed JS module with a REPL backend.
type API struct {
	Greeting string
	Run      func(string) (sobek.Value, error)
	Block    func()
}

// NewModuleInstance returns a new module instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	api := &API{
		Greeting: "WELCOME TO REPLIT!",
		Run: func(code string) (sobek.Value, error) {
			return vu.Runtime().RunString(code)
		},
	}
	api.Block = func() {
		startREPLServer(api)
	}

	return &ModuleInstance{
		module: &Module{API: api},
	}
}

// Start a TCP REPL server to accept and execute commands.
func startREPLServer(api *API) {
	ln, err := net.Listen("tcp", ":8089") // Listen on port 8089
	if err != nil {
		fmt.Println("Error starting REPL server:", err)
		return
	}
	fmt.Println("REPL server running on port 8089...")

	var wg sync.WaitGroup
	shutdown := make(chan struct{})
	once := new(sync.Once) // Ensures shutdown is only triggered once

	go func() {
		<-shutdown
		ln.Close()
		wg.Wait() // Ensure all goroutines finish before exiting
		fmt.Println("Server shutting down...")
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-shutdown:
				return // Stop accepting new connections on shutdown
			default:
				fmt.Println("Error accepting connection:", err)
			}
			continue
		}

		wg.Add(1)
		go func() {
			handleConnection(conn, api, shutdown, once)
			wg.Done()
		}()
	}
}

// Handle incoming REPL commands.
func handleConnection(conn net.Conn, api *API, shutdown chan struct{}, once *sync.Once) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("> ")) // REPL prompt
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		cmd = strings.TrimSpace(cmd)
		if cmd == "exit" {
			conn.Write([]byte("Exiting REPL...\n"))
			once.Do(func() { close(shutdown) }) // Close only once
			return
		}

		result, err := api.Run(cmd)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		} else {
			conn.Write([]byte(fmt.Sprintf("%v\n", result)))
		}
	}
}
