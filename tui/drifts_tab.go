package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
)

type DriftsTab struct {
	c            *client.Client
	err          error
	isLoading    bool
	loading      spinner.Model
	currentDrift client.Drift
	currentAsset client.Asset
	driftsVP     *ChooserViewport
	assetsVP     *ChooserViewport
	codifyVP     *Viewport
	width        int
	height       int
	ready        bool
}

func NewDriftsTab(c *client.Client) *DriftsTab {
	loading := spinner.NewModel()
	loading.Spinner = spinner.Dot
	loading.Style = purpleText

	return &DriftsTab{
		c:         c,
		loading:   loading,
		isLoading: true,
	}
}

func (envs *DriftsTab) Key() string                 { return "2" }
func (envs *DriftsTab) Name() string                { return "Inventory" }
func (envs *DriftsTab) NormalStyle() lipgloss.Style { return normalTabStyle }
func (envs *DriftsTab) ActiveStyle() lipgloss.Style { return activeTabStyle }

func (m *DriftsTab) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, m.loadDrifts)
}

func (m *DriftsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.loading, cmd = m.loading.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case KeyEnter:
			if m.assetsVP != nil {
				return m, m.codify
			} else if m.driftsVP != nil {
				return m, m.showDrift
			}
		case KeyEsc, KeyBackspace:
			m.err = nil
			m.isLoading = false

			if m.codifyVP != nil {
				m.codifyVP = nil
			} else if m.assetsVP != nil {
				m.assetsVP = nil
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case successMsg:
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil
	}

	switch {
	case m.codifyVP != nil:
		_, cmd = m.codifyVP.Update(msg)
	case m.assetsVP != nil:
		_, cmd = m.assetsVP.Update(msg)
	case m.driftsVP != nil:
		_, cmd = m.driftsVP.Update(msg)
	}

	return m, cmd
}

func (m *DriftsTab) View() string {
	if !m.ready {
		return initializing
	}

	var b strings.Builder

	switch {
	case m.err != nil:
		fmt.Fprintf(&b, "Error: %s\n", m.err)
	case m.isLoading:
		fmt.Fprintf(&b, "%s Loading...\n\n", m.loading.View())
	case m.codifyVP != nil:
		b.WriteString(m.codifyVP.View())
	case m.assetsVP != nil:
		b.WriteString(m.assetsVP.View())
	default:
		b.WriteString(m.driftsVP.View())
	}

	return b.String()
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

func (m *DriftsTab) loadDrifts() tea.Msg {
	m.isLoading = true

	drifts, err := m.c.ListDrifts(true, 90)
	if err != nil {
		return errMsg{err}
	}

	var (
		driftIDWidth   = 31
		driftDateWidth = m.width - driftIDWidth - horizontalMargins
	)

	choices := make([]Choice, len(drifts))
	for i, drift := range drifts {
		driftID := trim(strings.Replace(drift.ID, "Drifts/", "", 1), driftIDWidth)
		date := trim(drift.CreationDate().Format(time.RFC1123), driftDateWidth)

		choices[i] = Choice{
			Object: drift,
			Display: fmt.Sprintf(
				"%s %s",
				typeStyle.Render(lipgloss.PlaceHorizontal(driftIDWidth, lipgloss.Center, driftID)),
				lipgloss.PlaceHorizontal(driftDateWidth, lipgloss.Center, date),
			),
			ChosenDisplay: highlightedStyle.Render(fmt.Sprintf(
				"%s %s",
				lipgloss.PlaceHorizontal(driftIDWidth, lipgloss.Center, driftID),
				lipgloss.PlaceHorizontal(driftDateWidth, lipgloss.Center, date),
			)),
		}
	}

	m.driftsVP = NewChooserViewport(choices, m.width, m.height)
	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}

func (m *DriftsTab) showDrift() tea.Msg {
	m.isLoading = true

	m.currentDrift, _ = m.driftsVP.CurrentChoice().Object.(client.Drift)

	assets, err := m.c.ShowDrift(m.currentDrift.ID)
	if err != nil {
		return errMsg{err}
	}

	var (
		assetTypeWidth  = 20
		assetStateWidth = 11
		assetIDWidth    = m.width - assetTypeWidth - assetStateWidth - horizontalMargins
	)

	choices := make([]Choice, len(assets))
	for i, asset := range assets {
		assetStyle := managedStyle
		switch asset.State {
		case client.StateModified:
			assetStyle = modifiedStyle
		case client.StateUnmanaged:
			assetStyle = unmanagedStyle
		}

		assetType := trim(asset.Type, assetTypeWidth)
		assetID := trim(asset.ID, assetIDWidth)
		assetState := trim(string(asset.State), assetStateWidth)

		choices[i] = Choice{
			Object: asset,
			Display: fmt.Sprintf(
				"%s %s %s",
				typeStyle.Render(lipgloss.PlaceHorizontal(assetTypeWidth, lipgloss.Center, assetType)),
				lipgloss.PlaceHorizontal(assetIDWidth, lipgloss.Left, assetID),
				assetStyle.Render(lipgloss.PlaceHorizontal(assetStateWidth, lipgloss.Center, assetState)),
			),
			ChosenDisplay: highlightedStyle.Render(fmt.Sprintf(
				"%s %s %s",
				lipgloss.PlaceHorizontal(assetTypeWidth, lipgloss.Center, assetType),
				lipgloss.PlaceHorizontal(assetIDWidth, lipgloss.Left, assetID),
				lipgloss.PlaceHorizontal(assetStateWidth, lipgloss.Center, assetState),
			)),
		}
	}

	m.assetsVP = NewChooserViewport(choices, m.width, m.height)
	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}

func (m *DriftsTab) codify() tea.Msg {
	m.isLoading = true

	m.currentAsset, _ = m.assetsVP.CurrentChoice().Object.(client.Asset)

	code, err := m.c.Codify(
		m.currentAsset.Type,
		fmt.Sprintf("%s:%s", m.currentAsset.AccountID, m.currentAsset.ID),
	)
	if err != nil {
		return errMsg{err}
	}

	var b strings.Builder
	for _, line := range strings.Split(code, "\n") {
		fmt.Fprintln(&b, strings.ReplaceAll(line, "\t", "    "))
	}

	m.codifyVP = NewViewport(b.String(), m.width, m.height)
	m.loading.Finish()
	m.isLoading = false

	return successMsg("success")
}

func trim(str string, maxWidth int) string {
	if len(str) > maxWidth {
		return str[0:maxWidth-1] + "…"
	}

	return str
}
