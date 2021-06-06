package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/infralight/cli/client"
	"golang.org/x/term"
)

type Model struct {
	width     int
	signIn    *SignInModel
	tabs      []Tab
	activeTab Tab
}

func Start(c *client.Client, accessKey, secretKey string) error {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	homeTab := NewHomeTab(c, width)

	return tea.NewProgram(
		&Model{
			width:  width,
			signIn: NewSignIn(c, accessKey, secretKey),
			tabs: []Tab{
				homeTab,
				NewEnvsTab(c, width),
				NewDriftsTab(c, width),
			},
			activeTab: homeTab,
		},

		// Use the full size of the terminal in its "alternate screen buffer"
		tea.WithAltScreen(),

		// Also turn on mouse support so we can track the mouse wheel
		tea.WithMouseCellMotion(),
	).Start()
}

func (m *Model) Init() tea.Cmd {
	return m.signIn.Init()
}

func (m *Model) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	if !m.signIn.isSignedIn {
		_, cmd = m.signIn.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == KeyCtrlC {
			return m, tea.Quit
		}

		// did the user ask to switch tabs?
		for _, tab := range m.tabs {
			if tab.Key() == msg.String() {
				m.activeTab = tab
				return m, tab.Init()
			}
		}

		_, cmd = m.activeTab.Update(msg)
	case tea.WindowSizeMsg:
		// pass this message to all tabs
		for _, tab := range m.tabs {
			_, cmd = tab.Update(msg)
		}
	default:
		_, cmd = m.activeTab.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	var b strings.Builder

	if !m.signIn.isSignedIn {
		// User is not signed in, render the Sign In component
		fmt.Fprintf(&b, "%s\n\n", headerStyle.Render("Infralight CLI"))
		b.WriteString(m.signIn.View())
		return b.String()
	}

	// User is signed in, display main UI

	// Render main menu
	{
		menuItems := make([]string, len(m.tabs))
		for i, tab := range m.tabs {
			style := tab.NormalStyle()
			if m.activeTab == tab {
				style = tab.ActiveStyle()
			}

			menuItems[i] = style.Render(fmt.Sprintf("%s [%s]", tab.Name(), tab.Key()))
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			menuItems...,
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		fmt.Fprintf(&b, "%s", row)
	}

	b.WriteString("\n")
	b.WriteString(m.activeTab.View())

	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
