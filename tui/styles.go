package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	headerHeight      = 5
	footerHeight      = 2
	verticalMargins   = headerHeight + footerHeight
	horizontalMargins = 8
)

var (
	purpleColor = lipgloss.Color("#a200ff")
	whiteColor  = lipgloss.Color("#fff")
	greyColor   = lipgloss.Color("240")

	headerStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(purpleColor).
			Foreground(whiteColor).
			Bold(true).
			Underline(true)

	purpleText = lipgloss.NewStyle().Foreground(purpleColor)
	greyText   = lipgloss.NewStyle().Foreground(greyColor)
	noStyle    = lipgloss.NewStyle()

	focusedSubmitButton = "[ " + purpleText.Render("Submit") + " ]"
	blurredSubmitButton = "[ " + greyText.Render("Submit") + " ]"

	highlightedStyle = lipgloss.NewStyle().
				Bold(true).
				Background(purpleColor).
				Foreground(whiteColor)
)
