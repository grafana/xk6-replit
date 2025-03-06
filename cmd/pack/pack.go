package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	const k6Path = "./k6" // temporary output path for the built k6 binary

	cmd := exec.Command("xk6", "build", "--with", "github.com/grafana/xk6-replit=.")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	log.Println("Building custom k6 binary...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to build custom k6: %v\nBuild output:\n%s", err, out.String())
	}
	log.Println("Custom k6 build complete.")

	tmpDir, err := os.MkdirTemp("", "replit-gen-")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}

	// Create an assets folder inside the temporary directory.
	assetsDir := filepath.Join(tmpDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		log.Fatalf("Failed to create assets directory: %v", err)
	}

	// Copy the built k6 binary into the assets folder.
	embeddedK6Path := filepath.Join(assetsDir, "k6")
	k6Data, err := os.ReadFile(k6Path)
	if err != nil {
		log.Fatalf("Failed to read custom k6 binary: %v", err)
	}
	if err := os.WriteFile(embeddedK6Path, k6Data, 0755); err != nil {
		log.Fatalf("Failed to write k6 binary to assets: %v", err)
	}

	source := `package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed assets/k6
var embeddedK6 []byte

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <script.js>", os.Args[0])
		os.Exit(1)
	}
	scriptPath := os.Args[1]

	// Validate that the script imports the required module and uses replit.with()
	scriptData, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Error reading script file: %v", err)
	}
	if !bytes.Contains(scriptData, []byte(` + "`import { replit } from \"k6/x/replit\";`" + `)) {
		log.Fatalln("Script must import k6/x/replit to use the replit feature.")
	}
	if !bytes.Contains(scriptData, []byte("await replit.with(")) {
		log.Fatalln("Script must use replit.with() to start the REPLIT.")
	}

	// Write the embedded k6 binary to a temporary file.
	tmpDir := os.TempDir()
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
`
	sourcePath := filepath.Join(tmpDir, "replit_embedded.go")
	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		log.Fatalf("Failed to write generated source file: %v", err)
	}

	binDir := "bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		log.Fatalf("Failed to create bin directory: %v", err)
	}

	for _, target := range []struct {
		GOOS   string
		GOARCH string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
		{"windows", "arm64"},
	} {
		// Determine final executable name in the format "replit-<GOOS>-<GOARCH>".
		exe := fmt.Sprintf("replit-%s-%s", target.GOOS, target.GOARCH)
		if target.GOOS == "windows" {
			exe += ".exe"
		}
		cmd := exec.Command("go", "build", "-o", exe, sourcePath)
		cmd.Dir = tmpDir
		cmd.Env = append(os.Environ(), "GOOS="+target.GOOS, "GOARCH="+target.GOARCH)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to build final executable for %s/%s: %v\nOutput:\n%s",
				target.GOOS, target.GOARCH, err, string(out))
		}
		log.Printf("Final executable built successfully for %s/%s: %s", target.GOOS, target.GOARCH, exe)

		// Move the final executable from tempDir to the bin directory.
		fexe := filepath.Join(tmpDir, exe)
		fexeDir := filepath.Join(binDir, exe)
		// If destination exists and is a file, remove it.
		if info, err := os.Stat(fexeDir); err == nil && info.IsDir() {
			log.Fatalf("Destination %s is an existing directory; aborting", fexeDir)
		} else if err == nil {
			if err := os.Remove(fexeDir); err != nil {
				log.Fatalf("Failed to remove existing file %s: %v", fexeDir, err)
			}
		}
		if err := os.Rename(fexe, fexeDir); err != nil {
			log.Fatalf("Failed to move final executable to bin directory: %v", err)
		}
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Printf("Warning: could not delete temporary directory %s: %v", tmpDir, err)
	}
	if err := os.Remove(k6Path); err != nil {
		log.Printf("Warning: could not remove custom k6 binary %s: %v", k6Path, err)
	}
}
