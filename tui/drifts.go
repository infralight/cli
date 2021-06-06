package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
)

type DriftsTab struct {
	c            *client.Client
	width        int
	isLoading    bool
	loading      spinner.Model
	drifts       []client.Drift
	assets       []client.Asset
	chooser      *ChooserModel
	currentDrift string
	err          error
	vp           viewport.Model
	ready        bool
}

func NewDriftsTab(c *client.Client, width int) *DriftsTab {
	loading := spinner.NewModel()
	loading.Spinner = spinner.Dot
	loading.Style = purpleText

	return &DriftsTab{
		c:         c,
		width:     width,
		loading:   loading,
		isLoading: true,
	}
}

func (envs *DriftsTab) Key() string                 { return "2" }
func (envs *DriftsTab) Name() string                { return "Drifts" }
func (envs *DriftsTab) NormalStyle() lipgloss.Style { return normalTabStyle }
func (envs *DriftsTab) ActiveStyle() lipgloss.Style { return activeTabStyle }

func (m *DriftsTab) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, m.loadDrifts)
}

func (m *DriftsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	const (
		headerHeight = 5
		footerHeight = 2
	)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		_, cmd = m.loading.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			if m.currentDrift == "" {
				return m, tea.Batch(spinner.Tick, m.showDrift)
			}
		case KeyEsc, KeyBackspace:
			m.currentDrift = ""
			return m, nil
		}
	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight

		if !m.ready {
			m.vp = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - verticalMargins,
			}
			m.vp.YPosition = headerHeight
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - verticalMargins
		}
	case successMsg:
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	if m.currentDrift == "" {
		_, cmd = m.chooser.Update(msg)
	} else {
		m.vp, cmd = m.vp.Update(msg)
	}

	return m, cmd
}

func (m *DriftsTab) View() string {
	var b strings.Builder

	switch {
	case m.err != nil:
		fmt.Fprintf(&b, "Error: %s\n", m.err)
	case m.isLoading:
		fmt.Fprintf(&b, "%s Loading...\n\n", m.loading.View())
	case m.currentDrift != "":
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#fff")).
			Padding(1).
			Width(m.vp.Width - 2).
			Height(m.vp.Height)
		b.WriteString(style.Render(m.vp.View()))
	default:
		b.WriteString(m.chooser.View())
	}

	return b.String()
}

func (m *DriftsTab) loadDrifts() tea.Msg {
	m.isLoading = true

	var err error
	m.drifts, err = m.c.ListDrifts(true, 50)
	if err != nil {
		return errMsg{err}
	}

	if len(m.drifts) > 0 {
		driftChoices := make([]Choice, len(m.drifts))
		for i, drift := range m.drifts {
			driftChoices[i] = Choice{
				ID:   drift.ID,
				Name: drift.CreationDate().Format(time.RFC1123),
			}
		}

		m.chooser = NewChooser("Latest Drifts:", driftChoices)
	}

	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}

func (m *DriftsTab) showDrift() tea.Msg {
	m.isLoading = true

	var err error
	m.assets, err = m.c.ShowDrift(m.chooser.CurrentChoice().ID)
	if err != nil {
		return errMsg{err}
	}

	m.currentDrift = m.chooser.CurrentChoice().ID
	m.vp.SetContent(m.assetList())
	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}

var (
	typeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#f5f5f5")).
			Foreground(lipgloss.Color("#111"))

	managedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#14ff08"))

	unmanagedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff0004"))

	modifiedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffb300"))
)

func (m *DriftsTab) assetList() string {
	var b strings.Builder
	for _, asset := range m.assets {
		assetStyle := managedStyle
		switch asset.State {
		case client.StateModified:
			assetStyle = modifiedStyle
		case client.StateUnmanaged:
			assetStyle = unmanagedStyle
		}

		fmt.Fprintf(
			&b,
			"%s %s %s\n",
			typeStyle.Render(lipgloss.PlaceHorizontal(20, lipgloss.Center, asset.Type)),
			lipgloss.PlaceHorizontal(m.vp.Width-38, lipgloss.Left, asset.ID),
			assetStyle.Render(lipgloss.PlaceHorizontal(11, lipgloss.Center, string(asset.State))),
		)
	}

	return b.String()
}
