// geno — unified entry point for the geno ecosystem.
//
// Usage:
//
//	geno                       open the interactive TUI (Bubbletea agent)
//	geno tt   [args...]        delegate to the `tt` CLI (geno-tt)
//	geno vault [args...]       delegate to `geno-vault`
//	geno surf  [args...]       delegate to `surf` (geno-surf)
//	geno tools [args...]       delegate to `geno-tools`
//	geno pear  [args...]       delegate to `pear` (geno-pear)
//	geno version               print version
//
// Subcommands are passed through verbatim to their respective binaries.
// If a binary is not on PATH, geno prints a helpful install hint.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/42euge/geno-cli/internal/app"
	"github.com/42euge/geno-cli/internal/config"
	"github.com/42euge/geno-cli/internal/ollama"
)

const version = "0.1.0"

// dispatch maps geno subcommand names to their binary names.
var dispatch = map[string]string{
	"tt":    "tt",
	"vault": "geno-vault",
	"surf":  "surf",
	"tools": "geno-tools",
	"pear":  "pear",
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runTUI()
		return
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Printf("geno %s\n", version)
	case "help", "--help", "-h":
		printHelp()
	default:
		bin, ok := dispatch[args[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "geno: unknown subcommand %q\n\n", args[0])
			printHelp()
			os.Exit(1)
		}
		delegate(bin, args[1:])
	}
}

// delegate exec-replaces the current process with the target binary.
// This means the subprocess inherits stdin/stdout/stderr cleanly and
// interactive programs (like geno-vault gui) behave correctly.
func delegate(bin string, args []string) {
	path, err := exec.LookPath(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"geno: %q not found on PATH.\n  Install with: brew install geno\n", bin)
		os.Exit(1)
	}
	argv := append([]string{path}, args...)
	if err := syscall.Exec(path, argv, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "geno: exec %s: %v\n", bin, err)
		os.Exit(1)
	}
}

func runTUI() {
	cfg := config.Parse()
	client := ollama.NewClient(cfg.OllamaURL)
	m := app.New(client, cfg.Model, cfg.NoTools)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "geno: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`geno — agentic workspace orchestration

Usage:
  geno                  open the interactive agent TUI
  geno tt   [args]      iTerm2 + workspace orchestration
  geno vault [args]     registry sync, web GUI, daemon
  geno surf  [args]     Chromium agent-side orchestration
  geno tools [args]     geno skillset manager
  geno pear  [args]     shared watch library
  geno version          print version

Install:
  brew tap 42euge/geno && brew install geno
`)
}
