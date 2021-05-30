package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Choice struct {
	ID   string
	Name string
}

type ChooserModel struct {
	title   string
	choices []Choice
	cursor  int
}

func NewChooser(title string, choices []Choice) *ChooserModel {
	return &ChooserModel{
		title:   title,
		choices: choices,
	}
}

func (m *ChooserModel) Init() tea.Cmd {
	return nil
}

func (m *ChooserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *ChooserModel) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", m.title)

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		fmt.Fprintf(&b, "%s %s\n", cursor, choice.Name)
	}

	return b.String()
}
