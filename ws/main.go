package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func repl(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		err = c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			panic(err)
		}

		mt, message, err := c.ReadMessage()
		if err != nil {
			panic(err)
		} else if mt != websocket.TextMessage {
			panic("binary message received")
		}
		if len(message) > 0 {
			fmt.Printf("%s\n", message)
		}
	}
}

func main() {
	http.HandleFunc("/repl", repl)
	go func() {
		err := exec.Command("xk6", "run", "-q", "ws/replit.js").Run()
		if err != nil {
			panic(err)
		}
	}()

	http.ListenAndServe("localhost:8080", nil)
}
