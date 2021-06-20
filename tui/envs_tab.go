package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
)

type EnvsTab struct {
	isLoading bool
	loading   spinner.Model
	chooser   *ChooserViewport
	err       error
	c         *client.Client
	width     int
	height    int
	ready     bool
}

func NewEnvsTab(c *client.Client) *EnvsTab {
	loading := spinner.NewModel()
	loading.Spinner = spinner.Dot
	loading.Style = purpleText

	return &EnvsTab{
		c:         c,
		loading:   loading,
		isLoading: true,
	}
}

func (envs *EnvsTab) Key() string                 { return "1" }
func (envs *EnvsTab) Name() string                { return "Environments" }
func (envs *EnvsTab) NormalStyle() lipgloss.Style { return normalTabStyle }
func (envs *EnvsTab) ActiveStyle() lipgloss.Style { return activeTabStyle }

func (m *EnvsTab) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, m.loadEnvs)
}

func (m *EnvsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		_, cmd = m.loading.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - verticalMargins
		m.ready = true
	case successMsg:
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	if m.chooser != nil {
		_, cmd = m.chooser.Update(msg)
	}

	return m, cmd
}

func (m *EnvsTab) View() string {
	if !m.ready {
		return initializing
	}

	var b strings.Builder

	switch {
	case m.err != nil:
		fmt.Fprintf(&b, "Error: %s\n", m.err)
	case m.isLoading:
		fmt.Fprintf(&b, "%s Loading...\n\n", m.loading.View())
	default:
		b.WriteString(m.chooser.View())
	}

	return b.String()
}

func (m *EnvsTab) loadEnvs() tea.Msg {
	m.isLoading = true

	envs, err := m.c.ListEnvironments()
	if err != nil {
		return errMsg{err}
	}

	choices := make([]Choice, len(envs))
	for i, env := range envs {
		choices[i] = Choice{
			Object:        env,
			Display:       env.Name,
			ChosenDisplay: highlightedStyle.Render(env.Name),
		}
	}

	m.chooser = NewChooserViewport(choices, m.width, m.height)
	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}
