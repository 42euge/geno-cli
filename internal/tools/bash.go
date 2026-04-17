package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/42euge/geno-cli/internal/ollama"
)

type BashTool struct{}

type bashArgs struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

func (t *BashTool) Definition() ollama.Tool {
	return makeTool("bash", "Execute a shell command and return stdout+stderr.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{"type": "string", "description": "Shell command to execute"},
			"timeout": map[string]any{"type": "integer", "description": "Timeout in seconds (default 30)"},
		},
		"required": []string{"command"},
	})
}

func (t *BashTool) Execute(ctx context.Context, rawArgs json.RawMessage) (string, error) {
	var args bashArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return "", fmt.Errorf("parse args: %w", err)
	}

	timeout := time.Duration(args.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", args.Command)
	out, err := cmd.CombinedOutput()

	result := string(out)
	if len(result) > 30000 {
		result = result[:30000] + "\n... (truncated)"
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return result + "\n(command timed out)", nil
		}
		return fmt.Sprintf("%s\nExit code: %s", result, err), nil
	}

	// Trim trailing whitespace for cleaner output
	return strings.TrimRight(result, "\n\t "), nil
}
