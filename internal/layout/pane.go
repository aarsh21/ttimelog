package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
	Title   string
	View    func() string
	Focused bool
	Width   int
	Height  int
}

const (
	paneBorderBlurColor  = lipgloss.Color("#3b4261")
	paneBorderFocusColor = lipgloss.Color("#7aa2f7")
)

var (
	paneBorderBlur = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(paneBorderBlurColor)

	paneBorderFocus = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(paneBorderFocusColor)
)

func drawTopBorder(title string, width int, focus bool) string {
	lineStyle := lipgloss.NewStyle().Foreground(paneBorderBlurColor)
	if focus {
		lineStyle = lipgloss.NewStyle().Foreground(paneBorderFocusColor)
	}

	leftCorner := lineStyle.Render("╭─")
	label := lineStyle.Render(title)

	usedWidth := lipgloss.Width(leftCorner) + lipgloss.Width(label) + 1 // +1 for right corner
	rightLine := lineStyle.Render(strings.Repeat("─", max(width-usedWidth, 1)) + "╮")

	return leftCorner + label + rightLine
}

func (p Pane) Render() string {
	border := paneBorderBlur

	if p.Focused {
		border = paneBorderFocus
	}

	view := p.View()
	content := border.Width(p.Width).Height(lipgloss.Height(view)).Render(view)

	return lipgloss.JoinVertical(lipgloss.Left, drawTopBorder(p.Title, p.Width+2, p.Focused), content)
}
