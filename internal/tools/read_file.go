package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/42euge/geno-cli/internal/ollama"
)

type ReadFileTool struct{}

type readFileArgs struct {
	Path   string `json:"path"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (t *ReadFileTool) Definition() ollama.Tool {
	return makeTool("read_file", "Read the contents of a file at the given path. Returns numbered lines.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":   map[string]any{"type": "string", "description": "File path to read"},
			"offset": map[string]any{"type": "integer", "description": "Line number to start from (0-based, default 0)"},
			"limit":  map[string]any{"type": "integer", "description": "Max lines to read (default 200)"},
		},
		"required": []string{"path"},
	})
}

func (t *ReadFileTool) Execute(_ context.Context, rawArgs json.RawMessage) (string, error) {
	var args readFileArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return "", fmt.Errorf("parse args: %w", err)
	}

	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")

	offset := args.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= len(lines) {
		return "(empty — offset beyond end of file)", nil
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 200
	}

	end := offset + limit
	if end > len(lines) {
		end = len(lines)
	}

	var sb strings.Builder
	for i := offset; i < end; i++ {
		fmt.Fprintf(&sb, "%4d\t%s\n", i+1, lines[i])
	}

	result := sb.String()
	if len(result) > 30000 {
		result = result[:30000] + "\n... (truncated)"
	}
	return result, nil
}
