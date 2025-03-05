package replit

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

//go:embed replit.js
var replitJS string

// API is the exposed JS module with a REPL backend.
type API struct {
	// NEW API
	Read  func() (string, error)
	Log   func(msg string)
	Warn  func(msg string)
	Error func(msg string)

	// Read from JS code
	Repl sobek.Value `js:"with"`
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
	// setup autocomplete
	var hc historyAutoCompleter
	if rl.Config.HistoryFile != "" {
		if err := hc.loadFromFile(rl.Config.HistoryFile); err != nil {
			panic(err)
		}
	}
	rl.Config.AutoComplete = &hc

	api.Read = func() (string, error) {
		line, err := readMultiLine(rl, multilineOpts{
			promptSingleline: ">>> ",
			promptMultiline:  "... ",
			eolMarker:        ";",
		})
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
	api.Warn = func(msg string) {
		color.Yellow(msg)
	}
	api.Error = func(msg string) {
		color.Red(msg)
	}

	rt := vu.Runtime()
	if err := rt.GlobalObject().Set("replit", api); err != nil {
		panic(err)
	}
	rjs, err := rt.RunString(replitJS)
	if err != nil {
		panic(err)
	}
	api.Repl = rjs.ToObject(rt).Get("repl") // get the repl function in the script

	return api
}

type multilineOpts struct {
	promptSingleline string
	promptMultiline  string
	eolMarker        string // end of line marker (e.g. ";")
}

func readMultiLine(rl *readline.Instance, opts multilineOpts) (string, error) {
	prompt := opts.promptSingleline           // start with single-line prompt
	defer rl.SetPrompt(opts.promptSingleline) // reset prompt

	var lines []string
	for {
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(line) != 0 {
				continue // skip partially typed input
			}
			continue // skip empty lines
		}
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if strings.TrimSpace(line) == "" { // skip empty lines
			continue
		}
		lines = append(lines, line)
		if strings.HasSuffix(line, opts.eolMarker) { // saw the marker; stop accumulation
			break
		}

		prompt = opts.promptMultiline // switch to multiline prompt
	}

	return strings.Join(lines, "\n"), nil
}

// historyAutoCompleter implements the readline.Completer interface.
// It suggests completions from the history slice.
type historyAutoCompleter struct {
	history []string
}

func (hc *historyAutoCompleter) loadFromFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	hc.history = strings.Split(string(buf), "\n")
	return nil
}

// Do returns, for each candidate in history that starts with the current input (up to pos),
// the suffix beyond that shared prefix. The returned length (pos) indicates how many characters
// of the line are common and should be replaced.
func (hc *historyAutoCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Skip completing empty lines.
	if strings.TrimSpace(string(line)) == "" {
		return nil, 0
	}
	// Use the input up to pos as the prefix to complete.
	prefix := string(line[:pos])
	var suggestions [][]rune
	for _, entry := range hc.history {
		if strings.HasPrefix(entry, prefix) {
			// Append the remainder of the entry (the part after the prefix).
			remainder := entry[len(prefix):]
			suggestions = append(suggestions, []rune(remainder))
		}
	}
	// Return the suggestions and the length of the common prefix.
	return suggestions, pos
}
