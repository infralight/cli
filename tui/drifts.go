package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
)

type DriftsTab struct {
	c     *client.Client
	width int
}

func NewDriftsTab(c *client.Client, width int) *DriftsTab {
	return &DriftsTab{
		c:     c,
		width: width,
	}
}

func (envs *DriftsTab) Key() string                 { return "2" }
func (envs *DriftsTab) Name() string                { return "Drifts" }
func (envs *DriftsTab) NormalStyle() lipgloss.Style { return normalTabStyle }
func (envs *DriftsTab) ActiveStyle() lipgloss.Style { return activeTabStyle }

func (m *DriftsTab) Init() tea.Cmd {
	return nil
}

func (m *DriftsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *DriftsTab) View() string {
	return "Nothing here yet"
}
