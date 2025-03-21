package replit

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
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
	Read      func() (string, error)
	Log       func(msg string)
	Warn      func(msg string)
	Error     func(msg string)
	Highlight func(msg, lang string)

	// Read from JS code
	Repl sobek.Value `js:"run"`
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
	api.Highlight = func(msg string, lang string) {
		fmt.Println(highlight(msg, lang))
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
}

func readMultiLine(rl *readline.Instance, opts multilineOpts) (string, error) {
	prompt := opts.promptSingleline           // start with single-line prompt
	defer rl.SetPrompt(opts.promptSingleline) // reset prompt

	var lines []string
	for {
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(line) >= 0 {
				// user interrupted multiline input, return
				// whatever we've got so far
				break
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

		// FIXME: The compiling logic here kind of needs to match whatever replit.js tries
		// to run later (use or not of async, return, etc.). Ideally the logic should only
		// be in one place.

		// Try lines as an expression
		_, err = sobek.Compile("", "(async function(){return "+strings.Join(lines, "\n")+"}())", false)
		if err == nil {
			break
		}
		// Try lines as a statement
		_, err = sobek.Compile("", "(async function(){"+strings.Join(lines, "\n")+"}())", false)
		if err == nil {
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

// highlight formats code with ANSI escape codes.
func highlight(code, lang string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, code, lang, "terminal", "monokai")
	if err != nil {
		return code
	}
	return buf.String()
}
