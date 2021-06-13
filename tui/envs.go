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
	isLoading  bool
	loading    spinner.Model
	envs       []client.Environment
	chooser    *ChooserModel
	currentEnv int
	err        error
	c          *client.Client
	width      int
}

func NewEnvsTab(c *client.Client, width int) *EnvsTab {
	loading := spinner.NewModel()
	loading.Spinner = spinner.Dot
	loading.Style = purpleText

	return &EnvsTab{
		c:          c,
		width:      width,
		loading:    loading,
		isLoading:  true,
		currentEnv: -1,
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
	case successMsg:
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	if m.currentEnv == -1 {
		_, cmd = m.chooser.Update(msg)
	}

	return m, cmd
}

func (m *EnvsTab) View() string {
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
	var err error
	m.envs, err = m.c.ListEnvironments()
	if err != nil {
		return errMsg{err}
	}

	if len(m.envs) > 0 {
		envChoices := make([]Choice, len(m.envs))
		for i, env := range m.envs {
			envChoices[i] = Choice{
				ID:   env.ID,
				Name: env.Name,
			}
		}

		m.chooser = NewChooser("Available Environments:", envChoices)
	}

	m.isLoading = false

	return successMsg("success")
}
