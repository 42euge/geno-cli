package app

import (
	"context"

	"github.com/42euge/geno-cli/internal/agent"
	"github.com/42euge/geno-cli/internal/ollama"
	"github.com/42euge/geno-cli/internal/render"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateIdle state = iota
	stateStreaming
	stateToolCall
)

// Model is the root bubbletea model.
type Model struct {
	viewport viewport.Model
	textarea textarea.Model
	spinner  spinner.Model
	state    state
	loop     *agent.Loop
	model    string

	// Content displayed in viewport
	content string

	// Current streaming assistant response
	streamBuf string

	// Stream channel for current response
	streamCh <-chan agent.StreamMsg

	// Stats
	promptTokens int
	evalTokens   int

	// Dimensions
	width  int
	height int

	// Cancel context for current stream
	cancel context.CancelFunc
}

func New(client *ollama.Client, model string, noTools bool) Model {
	ta := textarea.New()
	ta.Placeholder = "Ask Geno anything..."
	ta.Focus()
	ta.SetHeight(3)
	ta.SetWidth(80)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false) // Enter sends, not newline

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	loop := agent.NewLoop(client, model, noTools)
	modeStr := "agentic"
	if !loop.ToolsActive() {
		modeStr = "chat"
	}
	welcome := "Welcome to geno-cli! Using model: " + model + " (" + modeStr + " mode)\nType a message and press Enter.\n"

	vp := viewport.New(80, 20)
	vp.SetContent(welcome)

	return Model{
		viewport: vp,
		textarea: ta,
		spinner:  sp,
		state:    stateIdle,
		loop:     loop,
		model:    model,
		content:  welcome,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

// streamNextMsg reads the next message from the stream channel.
func streamNextMsg(ch <-chan agent.StreamMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return streamDoneMsg{}
		}
		return streamMsg(msg)
	}
}

// Message types for bubbletea
type streamMsg agent.StreamMsg
type streamDoneMsg struct{}

func (m *Model) appendContent(s string) {
	m.content += s
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}

func (m *Model) appendRendered(md string) {
	rendered := render.Markdown(md)
	m.content += rendered
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()
}
