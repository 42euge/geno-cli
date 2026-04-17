package chat

import "github.com/42euge/geno-cli/internal/ollama"

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
		out[i] = ollama.Message{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}
	return out
}
