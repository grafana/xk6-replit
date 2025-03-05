package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed assets/replit.js
var replitJS []byte

// Where to store the custom k6 binary once built.
var customK6Path = filepath.Join(".", "k6")

func main() {
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

	// 2. Now run the custom k6 binary with replit.js.
	fmt.Printf("Welcome to REPLIT!\n")

	// write repl to a temp file
	if err := os.WriteFile("replit.js", replitJS, 0644); err != nil {
		log.Fatalf("Error writing replit.js: %v", err)
	}
	defer func() {
		if err := os.Remove("replit.js"); err != nil {
			log.Printf("Error removing replit.js: %v", err)
		}
	}()

	runCmd := exec.Command("./k6", "run", "-q", "--no-summary", "replit.js")
	runCmd.Stdin = os.Stdin   // so we can type commands interactively
	runCmd.Stdout = os.Stdout // send REPL output to our terminal
	runCmd.Stderr = os.Stderr // send errors directly to our terminal

	if err := runCmd.Run(); err != nil {
		log.Fatalf("Error running k6: %v", err)
	}
}
