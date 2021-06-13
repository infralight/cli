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
	code         string
	currentDrift string
	assetCursor  int
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

			return m, tea.Batch(spinner.Tick, m.codify)
		case KeyEsc, KeyBackspace:
			m.err = nil
			m.isLoading = false
			if m.code != "" {
				m.code = ""
			} else {
				m.currentDrift = ""
			}
			return m, nil
		case "up", "k":
			if m.currentDrift != "" && m.assetCursor > 0 {
				m.assetCursor--
				m.vp.SetContent(m.assetList())

				if m.vp.YOffset > 0 && m.assetCursor-m.vp.Height-1 > int(0.1*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "pgup", "b":
			if m.currentDrift != "" && m.assetCursor > 0 {
				if m.assetCursor >= m.vp.Height {
					m.assetCursor -= m.vp.Height
				} else {
					m.assetCursor = 0
				}

				m.vp.SetContent(m.assetList())

				if m.vp.YOffset > 0 && m.assetCursor-m.vp.Height-1 > int(0.1*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "down", "j":
			if m.currentDrift != "" && m.assetCursor < len(m.assets)-1 {
				m.assetCursor++
				m.vp.SetContent(m.assetList())
				if m.assetCursor < int(0.9*float64(m.vp.Height)) {
					return m, nil
				}
			}
		case "pgdown", " ", "f":
			if m.currentDrift != "" && m.assetCursor < len(m.assets)-1 {
				m.assetCursor += m.vp.Height
				m.vp.SetContent(m.assetList())
				if m.assetCursor < int(0.9*float64(m.vp.Height)) {
					return m, nil
				}
			}
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
		return m, nil
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
		if m.code != "" {
			m.vp.SetContent(m.code)
		}

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
	m.assetCursor = 0

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

	highlightedStyle = lipgloss.NewStyle().
				Bold(true).
				Background(purpleColor).
				Foreground(whiteColor)
)

func (m *DriftsTab) assetList() string {
	var b strings.Builder
	for i, asset := range m.assets {
		assetStyle := managedStyle
		switch asset.State {
		case client.StateModified:
			assetStyle = modifiedStyle
		case client.StateUnmanaged:
			assetStyle = unmanagedStyle
		}

		if i == m.assetCursor {
			line := fmt.Sprintf(
				"%s %s %s",
				lipgloss.PlaceHorizontal(30, lipgloss.Center, asset.Type),
				lipgloss.PlaceHorizontal(m.vp.Width-41-8, lipgloss.Left, asset.ID),
				lipgloss.PlaceHorizontal(11, lipgloss.Center, string(asset.State)),
			)

			fmt.Fprintln(&b, highlightedStyle.Render(line))
		} else {
			fmt.Fprintf(
				&b,
				"%s %s %s\n",
				typeStyle.Render(lipgloss.PlaceHorizontal(30, lipgloss.Center, asset.Type)),
				lipgloss.PlaceHorizontal(m.vp.Width-41-8, lipgloss.Left, asset.ID),
				assetStyle.Render(lipgloss.PlaceHorizontal(11, lipgloss.Center, string(asset.State))),
			)
		}
	}

	return b.String()
}

func (m *DriftsTab) codify() tea.Msg {
	m.isLoading = true

	asset := m.assets[m.assetCursor]

	var err error
	m.code, err = m.c.Codify(asset.Type, asset.ID)
	if err != nil {
		m.code = "error"
		return errMsg{err}
	}

	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}
