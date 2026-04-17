package app

import (
	"context"
	"fmt"

	"github.com/42euge/geno-cli/internal/agent"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case "esc":
			if m.state == stateStreaming || m.state == stateToolCall {
				if m.cancel != nil {
					m.cancel()
				}
				m.state = stateIdle
				m.appendContent("\n(cancelled)\n")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.state != stateIdle {
				return m, nil
			}
			input := m.textarea.Value()
			if input == "" {
				return m, nil
			}
			m.textarea.Reset()

			// Show user message
			m.appendContent(fmt.Sprintf("\n\033[1;36m> %s\033[0m\n\n", input))

			// Start streaming
			m.state = stateStreaming
			m.streamBuf = ""
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			m.streamCh = m.loop.Send(ctx, input)

			return m, streamNextMsg(m.streamCh)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 0
		inputHeight := 5 // textarea + borders
		statusHeight := 1
		vpHeight := m.height - headerHeight - inputHeight - statusHeight
		if vpHeight < 1 {
			vpHeight = 1
		}
		m.viewport.Width = m.width
		m.viewport.Height = vpHeight
		m.textarea.SetWidth(m.width - 2)
		m.viewport.SetContent(m.content)
		return m, nil

	case streamMsg:
		sm := agent.StreamMsg(msg)
		return m.handleStreamMsg(sm)

	case streamDoneMsg:
		// Channel closed unexpectedly
		if m.state != stateIdle {
			m.state = stateIdle
			m.appendContent("\n")
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update sub-components
	if m.state == stateIdle {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleStreamMsg(sm agent.StreamMsg) (tea.Model, tea.Cmd) {
	switch {
	case sm.Chunk != "":
		m.streamBuf += sm.Chunk
		m.appendContent(sm.Chunk)
		return m, streamNextMsg(m.streamCh)

	case sm.ToolCall != nil:
		m.state = stateToolCall
		m.appendContent(fmt.Sprintf("\n\033[1;33m[tool: %s]\033[0m\n", sm.ToolCall.Name))
		return m, streamNextMsg(m.streamCh)

	case sm.ToolDone != nil:
		result := sm.ToolDone.Result
		if len(result) > 500 {
			result = result[:500] + "..."
		}
		m.appendContent(fmt.Sprintf("\033[2m%s\033[0m\n\n", result))
		m.state = stateStreaming
		m.streamBuf = ""
		return m, streamNextMsg(m.streamCh)

	case sm.Done != nil:
		m.state = stateIdle
		m.promptTokens = sm.Done.PromptTokens
		m.evalTokens = sm.Done.EvalTokens
		m.appendContent("\n")
		return m, nil

	case sm.Error != nil:
		m.state = stateIdle
		m.appendContent(fmt.Sprintf("\n\033[1;31mError: %s\033[0m\n", sm.Error))
		return m, nil

	default:
		return m, streamNextMsg(m.streamCh)
	}
}
