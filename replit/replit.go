package replit

import (
	_ "embed"
	"os"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

//go:embed replit.js
var jsCode string

// API is the exposed JS module with a REPL backend.
type API struct {
	// NEW API
	Read  func() (string, error)
	Log   func(msg string)
	Error func(msg string)
	Repl  func(sobek.Value) (sobek.Value, error)
}

// NewAPI returns a new API instance.
func NewAPI(vu modules.VU) *API {
	api := &API{}

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
		if line == "clear" {
			readline.ClearScreen(os.Stdout)
		}
		return line, nil
	}
	api.Log = func(msg string) {
		color.Green(msg)
	}
	api.Error = func(msg string) {
		color.Red(msg)
	}

	jsValues, err := vu.Runtime().RunString(jsCode)
	if err != nil {
		panic(err)
	}
	replFn := jsValues.ToObject(vu.Runtime()).Get("repl")

	api.Repl = func(v sobek.Value) (sobek.Value, error) {
		result, ok := sobek.AssertFunction(replFn)
		if !ok {
			panic("replit.js is incorrectly defined")
		}

		return result(replFn, v)
	}

	return api
}
