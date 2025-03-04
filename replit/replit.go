package replit

import (
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

// Response struct for JSON serialization
type Response struct {
	Status  string `json:"status"`
	Output  string `json:"output,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// API is the exposed JS module with a REPL backend.
type API struct {
	// OLD API
	Greeting string
	Run      func(string) (sobek.Value, error)
	Block    func()

	// NEW API
	Read  func() (string, error)
	Log   func(msg string)
	Error func(msg string)
}

// NewAPI returns a new API instance.
func NewAPI(vu modules.VU) *API {
	api := &API{
		Greeting: "WELCOME TO REPLIT!",
		Run: func(code string) (sobek.Value, error) {
			return vu.Runtime().RunString(code)
		},
	}
	api.Block = func() {
		startREPLServer(api)
	}

	// NEW API

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/k6-repl-history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	api.Read = func() (string, error) {
		line, err := rl.Readline()
		if err != nil {
			return "", err
		}
		return line, nil
	}
	api.Log = func(msg string) {
		color.Green(msg)
	}
	api.Error = func(msg string) {
		color.Red(msg)
	}

	return api
}
