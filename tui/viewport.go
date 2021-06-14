package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Viewport struct {
	vp      viewport.Model
	content string
}

func NewViewport(content string, width, height int) *Viewport {
	m := &Viewport{
		content: content,
		vp: viewport.Model{
			Width:  width,
			Height: height - verticalMargins,
		},
	}
	m.vp.YPosition = headerHeight
	m.vp.SetContent(m.content)
	return m
}

func (m *Viewport) Init() tea.Cmd {
	return nil
}

func (m *Viewport) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *Viewport) View() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#fff")).
		Padding(1).
		Width(m.vp.Width - 2).
		Height(m.vp.Height)
	return style.Render(m.vp.View())
}
