package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/42euge/geno-cli/internal/ollama"
)

type GrepTool struct{}

type grepArgs struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
	Glob    string `json:"glob"`
}

func (t *GrepTool) Definition() ollama.Tool {
	return makeTool("grep", "Search file contents for a regex pattern using grep.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern": map[string]any{"type": "string", "description": "Regex pattern to search for"},
			"path":    map[string]any{"type": "string", "description": "Directory or file to search in (default: current dir)"},
			"glob":    map[string]any{"type": "string", "description": "File glob filter (e.g. '*.go')"},
		},
		"required": []string{"pattern"},
	})
}

func (t *GrepTool) Execute(ctx context.Context, rawArgs json.RawMessage) (string, error) {
	var args grepArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return "", fmt.Errorf("parse args: %w", err)
	}

	cmdArgs := []string{"-rn", "--color=never"}
	if args.Glob != "" {
		cmdArgs = append(cmdArgs, "--include="+args.Glob)
	}
	cmdArgs = append(cmdArgs, args.Pattern)

	searchPath := "."
	if args.Path != "" {
		searchPath = args.Path
	}
	cmdArgs = append(cmdArgs, searchPath)

	cmd := exec.CommandContext(ctx, "grep", cmdArgs...)
	out, err := cmd.CombinedOutput()

	result := string(out)
	if len(result) > 30000 {
		result = result[:30000] + "\n... (truncated)"
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "(no matches)", nil
		}
		return result, nil
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) > 100 {
		result = strings.Join(lines[:100], "\n") + fmt.Sprintf("\n... (%d more matches)", len(lines)-100)
	}

	return result, nil
}
