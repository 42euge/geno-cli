package render

import (
	"github.com/charmbracelet/glamour"
)

var renderer *glamour.TermRenderer

func init() {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// Fallback: no rendering
		return
	}
	renderer = r
}

func Markdown(s string) string {
	if renderer == nil || s == "" {
		return s
	}
	out, err := renderer.Render(s)
	if err != nil {
		return s
	}
	return out
}
