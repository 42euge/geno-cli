package ollama

import "encoding/json"

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools,omitempty"`
	Options  *Options  `json:"options,omitempty"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type ChatResponse struct {
	Model              string  `json:"model"`
	Message            Message `json:"message"`
	Done               bool    `json:"done"`
	DoneReason         string  `json:"done_reason,omitempty"`
	TotalDuration      int64   `json:"total_duration,omitempty"`
	PromptEvalCount    int     `json:"prompt_eval_count,omitempty"`
	EvalCount          int     `json:"eval_count,omitempty"`
	EvalDuration       int64   `json:"eval_duration,omitempty"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumCtx      int     `json:"num_ctx,omitempty"`
}
