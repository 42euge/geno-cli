package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	inputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
)

func (m Model) View() string {
	// Status bar
	stateStr := "ready"
	switch m.state {
	case stateStreaming:
		stateStr = m.spinner.View() + " streaming"
	case stateToolCall:
		stateStr = m.spinner.View() + " running tool"
	}

	tokens := ""
	if m.evalTokens > 0 {
		tokens = fmt.Sprintf(" | tokens: %d/%d", m.promptTokens, m.evalTokens)
	}

	modeStr := "agentic"
	if !m.loop.ToolsActive() {
		modeStr = "chat"
	}

	status := statusStyle.Width(m.width).Render(
		fmt.Sprintf(" geno-cli | %s (%s) | %s%s", m.model, modeStr, stateStr, tokens),
	)

	// Input area
	input := inputBorder.Width(m.width - 4).Render(m.textarea.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		status,
		input,
	)
}
