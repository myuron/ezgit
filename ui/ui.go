package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/myuron/ezgit/gitcontroller"
)

// Model ... A structure for managing the application's state
type Model struct {
	prefixes      []string
	cursor        int
	prefix        string
	currentScreen screen
	textInput     textinput.Model
	Err           error
}

type screen int

const (
	prefixScreen screen = iota
	messageScreen
)

// Init ... A function that returns an initial command for the application to run
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update ... A function that handles incoming events and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case prefixScreen:
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "j":
				if m.cursor < len(m.prefixes)-1 {
					m.cursor++
				}
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "enter":
				m.prefix = m.prefixes[m.cursor]
				m.currentScreen = messageScreen
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
	case messageScreen:
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter":
				if strings.TrimSpace(m.textInput.Value()) == "" {
					return m, nil
				}
				s := m.prefix + ": " + m.textInput.Value()
				if err := gitcontroller.Commit(s); err != nil {
					m.Err = err
				}
				return m, tea.Quit
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View ... A function that renders the UI based on the data in the model
func (m Model) View() tea.View {
	switch m.currentScreen {
	case prefixScreen:
		s := "Select Prefix\n\n"
		for i, choice := range m.prefixes {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		return tea.NewView(s)
	case messageScreen:
		s := "Write Message\n\n"
		s += m.prefix
		s += "\n"
		s += m.textInput.View()
		return tea.NewView(s)
	}
	return tea.NewView("")
}

// InitialModel ... Restore to default settings
func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "description..."
	ti.Focus()
	ti.CharLimit = 72 // Maximum number of characters in a commit title before it wraps
	ti.SetWidth(72)
	return Model{
		prefixes: []string{
			"fix",
			"feat",
			"build",
			"chore",
			"ci",
			"docs",
			"perf",
			"refactor",
			"revert",
			"style",
			"test",
		},
		prefix:        "",
		currentScreen: prefixScreen,
		textInput:     ti,
	}
}
