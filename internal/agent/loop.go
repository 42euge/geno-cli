package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/42euge/geno-cli/internal/chat"
	"github.com/42euge/geno-cli/internal/ollama"
	"github.com/42euge/geno-cli/internal/tools"
)

const MaxToolRounds = 10

type Loop struct {
	client        *ollama.Client
	model         string
	history       *chat.History
	registry      *tools.Registry
	noTools       bool
	toolsDisabled bool // set true if model doesn't support tools
}

func NewLoop(client *ollama.Client, model string, noTools bool) *Loop {
	l := &Loop{
		client:   client,
		model:    model,
		history:  chat.NewHistory(),
		registry: tools.NewRegistry(),
		noTools:  noTools,
	}
	l.history.Add(chat.Message{Role: chat.RoleSystem, Content: SystemPrompt})
	return l
}

// StreamMsg represents a message from the agent loop to the TUI.
type StreamMsg struct {
	// Exactly one of these is set per message
	Chunk    string           // streaming text chunk
	ToolCall *ToolCallInfo    // model wants to call a tool
	ToolDone *ToolResultInfo  // tool execution result
	Done     *DoneInfo        // stream complete
	Error    error            // something went wrong
}

type ToolCallInfo struct {
	Name string
	Args string
}

type ToolResultInfo struct {
	Name   string
	Result string
}

type DoneInfo struct {
	PromptTokens int
	EvalTokens   int
}

// Send adds a user message and starts streaming the response.
// The returned channel emits StreamMsg until closed.
func (l *Loop) Send(ctx context.Context, userInput string) <-chan StreamMsg {
	out := make(chan StreamMsg, 64)
	l.history.Add(chat.Message{Role: chat.RoleUser, Content: userInput})

	go l.run(ctx, out)
	return out
}

func (l *Loop) run(ctx context.Context, out chan<- StreamMsg) {
	defer close(out)

	for round := 0; round < MaxToolRounds; round++ {
		req := ollama.ChatRequest{
			Model:    l.model,
			Messages: l.history.ToOllama(),
			Stream:   true,
		}
		if !l.noTools && !l.toolsDisabled {
			req.Tools = l.registry.Definitions()
		}

		ch, err := l.client.Chat(ctx, req)
		if err != nil {
			// If model doesn't support tools, retry without them
			if !l.toolsDisabled && strings.Contains(err.Error(), "does not support tools") {
				l.toolsDisabled = true
				req.Tools = nil
				ch, err = l.client.Chat(ctx, req)
				if err != nil {
					out <- StreamMsg{Error: err}
					return
				}
			} else {
				out <- StreamMsg{Error: err}
				return
			}
		}

		var fullContent string
		var toolCalls []ollama.ToolCall
		var doneInfo DoneInfo

		for resp := range ch {
			if resp.Message.Content != "" {
				fullContent += resp.Message.Content
				out <- StreamMsg{Chunk: resp.Message.Content}
			}
			if len(resp.Message.ToolCalls) > 0 {
				toolCalls = append(toolCalls, resp.Message.ToolCalls...)
			}
			if resp.Done {
				doneInfo.PromptTokens = resp.PromptEvalCount
				doneInfo.EvalTokens = resp.EvalCount
			}
		}

		// Add assistant message to history (including tool calls if any)
		assistantMsg := chat.Message{Role: chat.RoleAssistant, Content: fullContent}
		for _, tc := range toolCalls {
			assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, chat.ToolCallRecord{
				Name:      tc.Function.Name,
				Arguments: string(tc.Function.Arguments),
			})
		}
		l.history.Add(assistantMsg)

		if len(toolCalls) == 0 {
			out <- StreamMsg{Done: &doneInfo}
			return
		}

		// Execute tool calls
		for _, tc := range toolCalls {
			argsStr := string(tc.Function.Arguments)
			out <- StreamMsg{ToolCall: &ToolCallInfo{Name: tc.Function.Name, Args: argsStr}}

			result, err := l.registry.Execute(ctx, tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				result = fmt.Sprintf("Error: %s", err)
			}

			out <- StreamMsg{ToolDone: &ToolResultInfo{Name: tc.Function.Name, Result: result}}

			// Add tool result to history
			l.history.Add(chat.Message{
				Role:     chat.RoleTool,
				Content:  result,
				ToolName: tc.Function.Name,
			})
		}
		// Loop continues — send history back to model with tool results
	}

	out <- StreamMsg{Error: fmt.Errorf("exceeded maximum tool rounds (%d)", MaxToolRounds)}
}

// ToolResultMessage creates an Ollama message for a tool result.
func ToolResultMessage(name, content string) ollama.Message {
	return ollama.Message{
		Role:    "tool",
		Content: content,
	}
}

// ToolsActive returns true if the model is using tools (not in chat-only mode).
func (l *Loop) ToolsActive() bool {
	return !l.noTools && !l.toolsDisabled
}

// History returns the raw conversation history (exported for serialization).
func (l *Loop) History() []chat.Message {
	return l.history.Messages
}

// FormatToolArgs formats raw JSON arguments for display.
func FormatToolArgs(raw json.RawMessage) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return string(raw)
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
