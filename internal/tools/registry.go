package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/42euge/geno-cli/internal/ollama"
)

type Handler interface {
	Definition() ollama.Tool
	Execute(ctx context.Context, args json.RawMessage) (string, error)
}

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	r := &Registry{handlers: make(map[string]Handler)}
	r.Register(&ReadFileTool{})
	r.Register(&WriteFileTool{})
	r.Register(&ListFilesTool{})
	r.Register(&BashTool{})
	r.Register(&GrepTool{})
	return r
}

func (r *Registry) Register(h Handler) {
	def := h.Definition()
	r.handlers[def.Function.Name] = h
}

func (r *Registry) Definitions() []ollama.Tool {
	defs := make([]ollama.Tool, 0, len(r.handlers))
	for _, h := range r.handlers {
		defs = append(defs, h.Definition())
	}
	return defs
}

func (r *Registry) Execute(ctx context.Context, name string, args json.RawMessage) (string, error) {
	h, ok := r.handlers[name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", name)
	}
	return h.Execute(ctx, args)
}

func makeTool(name, description string, params any) ollama.Tool {
	p, _ := json.Marshal(params)
	return ollama.Tool{
		Type: "function",
		Function: ollama.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  p,
		},
	}
}
