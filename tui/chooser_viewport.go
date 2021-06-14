package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Choice struct {
	Object        interface{}
	Display       string
	ChosenDisplay string
}

type ChooserViewport struct {
	vp      viewport.Model
	cursor  int
	choices []Choice
}

func NewChooserViewport(
	choices []Choice,
	width int,
	height int,
) *ChooserViewport {
	m := &ChooserViewport{
		choices: choices,
		vp: viewport.Model{
			Width:  width,
			Height: height - verticalMargins,
		},
	}
	m.vp.YPosition = headerHeight
	return m
}

func (m *ChooserViewport) Init() tea.Cmd {
	return nil
}

func (m *ChooserViewport) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--

				if m.vp.YOffset > 0 && m.cursor-m.vp.Height-1 > int(0.1*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "pgup", "b":
			if m.cursor > 0 {
				if m.cursor >= m.vp.Height {
					m.cursor -= m.vp.Height
				} else {
					m.cursor = 0
				}

				if m.vp.YOffset > 0 && m.cursor-m.vp.Height-1 > int(0.1*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				if m.cursor < int(0.9*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "pgdown", " ", "f":
			if m.cursor < len(m.choices)-1 {
				m.cursor += m.vp.Height
				if m.cursor < int(0.9*float64(m.vp.Height)) {
					return m, nil
				}
			}
		}
	}

	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *ChooserViewport) View() string {
	m.vp.SetContent(m.render())

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#fff")).
		Padding(1).
		Width(m.vp.Width - 2).
		Height(m.vp.Height)
	return style.Render(m.vp.View())
}

func (m *ChooserViewport) render() string {
	var b strings.Builder
	for i, choice := range m.choices {
		if i == m.cursor {
			fmt.Fprintln(&b, choice.ChosenDisplay)
		} else {
			fmt.Fprintln(&b, choice.Display)
		}
	}
	return b.String()
}

func (m *ChooserViewport) CurrentChoice() Choice {
	return m.choices[m.cursor]
}
