package tui

import (
	"github.com/charmbracelet/lipgloss"
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
)
