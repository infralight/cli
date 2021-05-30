package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/infralight/cli/client"
)

type statusMsg int

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

type SignInModel struct {
	index        int
	inputs       []textinput.Model
	submitButton string
	loading      spinner.Model
	isLoading    bool
	isSignedIn   bool
	err          error
	c            *client.Client
}

func NewSignIn(c *client.Client, accessKey, secretKey string) *SignInModel {
	accessKeyInput := textinput.NewModel()
	accessKeyInput.Placeholder = "Access Key"
	accessKeyInput.Focus()
	accessKeyInput.PromptStyle = purpleText
	accessKeyInput.TextStyle = purpleText

	secretKeyInput := textinput.NewModel()
	secretKeyInput.Placeholder = "Secret Key"
	secretKeyInput.EchoMode = textinput.EchoPassword
	secretKeyInput.EchoCharacter = '•'

	if accessKey != "" {
		accessKeyInput.SetValue(accessKey)

		if secretKey != "" {
			secretKeyInput.SetValue(secretKey)
		}
	}

	loading := spinner.NewModel()
	loading.Spinner = spinner.Dot
	loading.Style = purpleText

	return &SignInModel{
		c:            c,
		index:        0,
		inputs:       []textinput.Model{accessKeyInput, secretKeyInput},
		submitButton: blurredSubmitButton,
		loading:      loading,
		isLoading:    accessKey != "" && secretKey != "",
	}
}

func (m *SignInModel) Init() tea.Cmd {
	cmds := []tea.Cmd{textinput.Blink, spinner.Tick}
	if m.inputs[0].Value() != "" && m.inputs[1].Value() != "" {
		cmds = append(cmds, m.signIn)
	}

	return tea.Batch(cmds...)
}

func (m *SignInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.loading, cmd = m.loading.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case KeyTab, KeyShiftTab, KeyEnter, KeyUp, KeyDown:
			// Cycle between inputs
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == KeyEnter && m.index == len(m.inputs) {
				// Sign-in to the Infralight App Server
				m.isLoading = true
				return m, m.signIn
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
	case statusMsg:
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds[i] = cmd
	}

	return m, tea.Batch(cmds...)
}

func (m *SignInModel) View() string {
	var b strings.Builder

	switch {
	case m.err != nil:
		fmt.Fprintf(&b, "Error: %s\n", m.err)
	case m.isLoading:
		fmt.Fprintf(&b, "%s Signing in...\n\n", m.loading.View())
	case !m.isSignedIn:
		for i := 0; i < len(m.inputs); i++ {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteString("\n")
			}
		}

		fmt.Fprintf(&b, "\n\n%s\n", m.submitButton)
	}

	return b.String()
}

func (m *SignInModel) signIn() tea.Msg {
	err := m.c.Authenticate(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return errMsg{err}
	}

	m.isLoading = false
	m.isSignedIn = true

	return statusMsg(200)
}
