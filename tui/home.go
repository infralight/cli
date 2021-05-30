package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
)

var (
	homeTabStyle = normalTabStyle.Copy().
		Foreground(purpleColor).
		Border(headerTabBorder, true).
		Bold(true)
)

type HomeTab struct {
	width int
}

func NewHomeTab(_ *client.Client, width int) *HomeTab {
	return &HomeTab{width}
}

func (home *HomeTab) Key() string                 { return "0" }
func (home *HomeTab) Name() string                { return "Infralight" }
func (home *HomeTab) NormalStyle() lipgloss.Style { return homeTabStyle }
func (home *HomeTab) ActiveStyle() lipgloss.Style { return homeTabStyle }

func (home *HomeTab) Init() tea.Cmd {
	return nil
}

func (m *HomeTab) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	return m, nil
}

func (m *HomeTab) View() string {
	return fmt.Sprintf(
		"%s\n\n",
		lipgloss.Place(
			m.width,
			15,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.
				NewStyle().
				Width(50).
				Height(2).
				Align(lipgloss.Center).
				Render(`Welcome to the Infralight CLI!
Select a tab from the above menu by
pressing the appropriate keyboard shortcut`),
			lipgloss.WithWhitespaceChars("※"),
			lipgloss.WithWhitespaceForeground(subtle),
		),
	)
}
