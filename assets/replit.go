package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed k6
var embeddedK6 []byte

// TODO:
// We want to import more modules by default (maybe even all of them?)
// or allow the user to somehow specify modules via CLI args.
const standalonejs = `
import http from "k6/http";
import { browser } from "k6/browser";
import { replit } from "k6/x/replit";
import { sleep } from "k6";

export const options = {
    scenarios: {
        ui: {
            executor: "shared-iterations",
            options: {
                browser: {
                    type: "chromium",
                },
            },
        },
    },
};

export default async function () {
    await replit.with({http: http, browser: browser, sleep: sleep});
}
`

func main() {
	tmpDir := os.TempDir()
	var scriptPath string

	if len(os.Args) < 2 {
		// If no script was specified, then run a default one which just
		// runs a standalone REPL.
		scriptPath = filepath.Join(tmpDir, "repl.js")
		if err := os.WriteFile(scriptPath, []byte(standalonejs), 0644); err != nil {
			log.Fatalf("Error writing embedded repl.js: %v", err)
		}
	} else {
		scriptPath = os.Args[1]
	}

	// Write the embedded k6 binary to a temporary file.
	customK6Path := filepath.Join(tmpDir, "custom_k6")
	if err := os.WriteFile(customK6Path, embeddedK6, 0755); err != nil {
		log.Fatalf("Error writing embedded k6 binary: %v", err)
	}

	fmt.Println("Welcome to REPLIT!")
	// Execute the custom k6 binary with the provided script.
	runCmd := exec.Command(customK6Path, "run", "-q", "--no-summary", scriptPath)
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		log.Fatalf("Error running k6: %v", err)
	}
}
