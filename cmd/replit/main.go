package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "embed"
)

// Where to store the custom k6 binary once built.
var customK6Path = filepath.Join(".", "k6")

func main() {
	var scriptPath string
	if len(os.Args) < 2 {
		exit(1, "Usage: replit <script.js>")
	}
	scriptPath = os.Args[1]

	// check if the script imports k6/replit
	scriptFile, err := os.Open(scriptPath)
	if err != nil {
		exit(1, fmt.Sprintf("Error opening script file: %v", err))
	}
	defer scriptFile.Close()

	buf, err := io.ReadAll(scriptFile)
	if err != nil {
		exit(1, fmt.Sprintf("Error reading script file: %v", err))
	}
	if !bytes.Contains(buf, []byte(`import { replit } from "k6/x/replit";`)) {
		exit(1, strings.TrimSpace(`
The script must import k6/x/replit to use the replit feature.

Usage:
	import { replit } from "k6/x/replit";
	// other imports...
	// your script code...`))
	}
	if !bytes.Contains(buf, []byte(`await replit.with(`)) {
		exit(1, strings.TrimSpace(`
The script must use replit.with() to start the REPL.

Usage:
	await replit.with();
	// or:
	await replit.with({myVariable: 42});`))
	}

	// 1. Check if the custom k6 binary already exists.
	if _, err := os.Stat(customK6Path); os.IsNotExist(err) {
		log.Println("Building custom k6 binary...")

		// We'll capture xk6's output, but won't show it unless there's an error.
		var out bytes.Buffer
		buildCmd := exec.Command("xk6", "build",
			"--with", "github.com/grafana/xk6-replit=.",
		)
		// Capture stdout/stderr in the same buffer:
		buildCmd.Stdout = &out
		buildCmd.Stderr = &out

		if err := buildCmd.Run(); err != nil {
			// Show the output only if there's an error.
			log.Fatalf("Failed to build custom k6: %v\n\nBuild output:\n%s", err, out.String())
		}

		log.Println("Custom k6 build complete.")
	}

	// 2. Now run the custom k6 binary with repl.js.
	fmt.Printf("Welcome to REPLIT!\n")

	runCmd := exec.Command("./k6", "run", "-q", "--no-summary", scriptPath)
	runCmd.Stdin = os.Stdin   // so we can type commands interactively
	runCmd.Stdout = os.Stdout // send REPL output to our terminal
	runCmd.Stderr = os.Stderr // send errors directly to our terminal

	if err := runCmd.Run(); err != nil {
		log.Fatalf("Error running k6: %v", err)
	}
}

func exit(status int, message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(status)
}
