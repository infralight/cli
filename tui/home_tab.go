package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
	"github.com/infralight/cli/version"
)

var (
	homeTabStyle = normalTabStyle.Copy().
		Foreground(purpleColor).
		Border(headerTabBorder, true).
		Bold(true)
)

type HomeTab struct {
	width  int
	height int
	ready  bool
}

func NewHomeTab(_ *client.Client) *HomeTab {
	return &HomeTab{}
}

func (m *HomeTab) Key() string                 { return "0" }
func (m *HomeTab) Name() string                { return version.Product }
func (m *HomeTab) NormalStyle() lipgloss.Style { return homeTabStyle }
func (m *HomeTab) ActiveStyle() lipgloss.Style { return homeTabStyle }

func (m *HomeTab) Init() tea.Cmd {
	return nil
}

func (m *HomeTab) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height - verticalMargins
		m.ready = true
	}
	return m, nil
}

func (m *HomeTab) View() string {
	if !m.ready {
		return initializing
	}

	return fmt.Sprintf(
		"%s\n\n",
		lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.
				NewStyle().
				Width(50).
				Height(3).
				Align(lipgloss.Center).
				Render(fmt.Sprintf(`Welcome to the %s CLI!
Select a tab from the above menu by
pressing the appropriate keyboard shortcut`, version.Product)),
			lipgloss.WithWhitespaceChars("※"),
			lipgloss.WithWhitespaceForeground(subtle),
		),
	)
}
