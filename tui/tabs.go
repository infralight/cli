package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab interface {
	tea.Model
	Key() string
	Name() string
	NormalStyle() lipgloss.Style
	ActiveStyle() lipgloss.Style
}

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	headerTabBorder = lipgloss.Border{
		Bottom:      "─",
		BottomLeft:  "──",
		BottomRight: "──",
	}

	normalTabStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 1)

	activeTabStyle = normalTabStyle.Copy().Border(activeTabBorder, true)

	tabGap = normalTabStyle.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)
