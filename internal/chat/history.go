package chat

import (
	"encoding/json"

	"github.com/42euge/geno-cli/internal/ollama"
)

type History struct {
	Messages []Message
}

func NewHistory() *History {
	return &History{}
}

func (h *History) Add(msg Message) {
	h.Messages = append(h.Messages, msg)
}

func (h *History) ToOllama() []ollama.Message {
	out := make([]ollama.Message, len(h.Messages))
	for i, m := range h.Messages {
		msg := ollama.Message{
			Role:    string(m.Role),
			Content: m.Content,
		}
		// Preserve tool calls on assistant messages
		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, ollama.ToolCall{
					Function: ollama.ToolCallFunction{
						Name:      tc.Name,
						Arguments: json.RawMessage(tc.Arguments),
					},
				})
			}
		}
		out[i] = msg
	}
	return out
}
