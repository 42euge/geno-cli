package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/42euge/geno-cli/internal/ollama"
)

type ListFilesTool struct{}

type listFilesArgs struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
}

func (t *ListFilesTool) Definition() ollama.Tool {
	return makeTool("list_files", "List files matching a glob pattern in a directory.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern": map[string]any{"type": "string", "description": "Glob pattern (e.g. '*.go', '**/*.ts')"},
			"path":    map[string]any{"type": "string", "description": "Directory to search in (default: current dir)"},
		},
		"required": []string{"pattern"},
	})
}

func (t *ListFilesTool) Execute(_ context.Context, rawArgs json.RawMessage) (string, error) {
	var args listFilesArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return "", fmt.Errorf("parse args: %w", err)
	}

	pattern := args.Pattern
	if args.Path != "" {
		pattern = filepath.Join(args.Path, pattern)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("glob: %w", err)
	}

	if len(matches) == 0 {
		return "(no matches)", nil
	}

	if len(matches) > 200 {
		matches = matches[:200]
	}

	return strings.Join(matches, "\n"), nil
}
