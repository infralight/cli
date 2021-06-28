package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/infralight/cli/client"
	"github.com/infralight/cli/config"
)

type ConfigureModel struct {
	index          int
	inputs         []textinput.Model
	message        string
	submitButton   string
	done           bool
	createdProfile string
}

func StartConfigure(showMessage string) (profile string, err error) {
	profileInput := textinput.NewModel()
	profileInput.Placeholder = "Profile [default]"
	profileInput.PromptStyle = purpleText
	profileInput.TextStyle = purpleText
	profileInput.Focus()

	accessKeyInput := textinput.NewModel()
	accessKeyInput.Placeholder = "Access Key"
	accessKeyInput.PromptStyle = purpleText
	accessKeyInput.TextStyle = purpleText

	secretKeyInput := textinput.NewModel()
	secretKeyInput.Placeholder = "Secret Key"
	secretKeyInput.EchoMode = textinput.EchoPassword
	secretKeyInput.EchoCharacter = '•'

	urlInput := textinput.NewModel()
	urlInput.Placeholder = fmt.Sprintf(
		"Infralight URL [%s]",
		client.DefaultInfralightURL,
	)
	urlInput.PromptStyle = purpleText
	urlInput.TextStyle = purpleText

	authHeaderInput := textinput.NewModel()
	authHeaderInput.Placeholder = fmt.Sprintf(
		"Authorization Header [%s]",
		client.DefaultAuthHeader,
	)
	authHeaderInput.PromptStyle = purpleText
	authHeaderInput.TextStyle = purpleText

	m := &ConfigureModel{
		message: showMessage,
		inputs: []textinput.Model{
			profileInput,
			urlInput,
			authHeaderInput,
			accessKeyInput,
			secretKeyInput,
		},
		submitButton: blurredSubmitButton,
	}

	err = tea.NewProgram(m).Start()
	if err != nil {
		return profile, err
	}

	return m.createdProfile, nil
}

func (m *ConfigureModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ConfigureModel) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyCtrlC:
			return m, tea.Quit
		case KeyTab, KeyShiftTab, KeyEnter, KeyUp, KeyDown:
			// Cycle between inputs
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == KeyEnter && m.index == len(m.inputs) {
				return m, m.writeConfig
			}

			// Cycle indexes
			if s == KeyUp || s == KeyShiftTab {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(m.inputs) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(m.inputs)
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.index {
					// Set focused state
					m.inputs[i].Focus()
					m.inputs[i].PromptStyle = purpleText
					m.inputs[i].TextStyle = purpleText
					continue
				}

				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			if m.index == len(m.inputs) {
				m.submitButton = focusedSubmitButton
			} else {
				m.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	case successMsg:
		m.done = true
		fmt.Printf("Configuration file saved to %s\n", string(msg))
		return m, tea.Quit
	case errMsg:
		m.done = true
		fmt.Printf("Error: %s\n", msg)
		return m, tea.Quit
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds[i] = cmd
	}

	return m, tea.Batch(cmds...)
}

func (m *ConfigureModel) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	if m.message != "" {
		fmt.Fprintf(&b, "%s:\n\n", m.message)
	}

	for i := 0; i < len(m.inputs); i++ {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	fmt.Fprintf(&b, "\n\n%s\n", m.submitButton)

	return b.String()
}

func (m *ConfigureModel) writeConfig() tea.Msg {
	c := config.Config{
		Profile:             m.inputs[0].Value(),
		URL:                 m.inputs[1].Value(),
		AuthorizationHeader: m.inputs[2].Value(),
		AccessKey:           m.inputs[3].Value(),
		SecretKey:           m.inputs[4].Value(),
	}

	if c.Profile == "" {
		c.Profile = "default"
	}
	if c.URL == "" {
		c.URL = client.DefaultInfralightURL
	}
	if c.AuthorizationHeader == "" {
		c.AuthorizationHeader = client.DefaultAuthHeader
	}
	if c.AccessKey == "" {
		return errMsg{errors.New("Access key must be provided")}
	}
	if c.SecretKey == "" {
		return errMsg{errors.New("Secret key must be provided")}
	}

	path, err := c.Save()
	if err != nil {
		return errMsg{err}
	}

	m.createdProfile = c.Profile

	return successMsg(path)
}
