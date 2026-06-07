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
	prefixDetails []PrefixDetail
	cursor        int
	choicedPrefix string 
	currentScreen screen
	textInput     textinput.Model
	Err           error
}

// PrefixDetail ... A structure for conbining prefix and description
type PrefixDetail struct {
	prefix      string
	description string
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
				if m.cursor < len(m.prefixDetails)-1 {
					m.cursor++
				}
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "enter":
				m.choicedPrefix = m.prefixDetails[m.cursor].prefix
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
				repoPath := "."
				commitMsg := m.choicedPrefix + ": " + m.textInput.Value()
				repo, err := gitcontroller.OpenRepository(repoPath)
				if err != nil {
					m.Err = err
					return m, tea.Quit
				}
				author, err := gitcontroller.LoadAuthor(repo)
				if err != nil {
					m.Err = err
					return m, tea.Quit
				}
				if err := gitcontroller.Commit(repo, commitMsg, author); err != nil {
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
		for i, p := range m.prefixDetails {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s: %s\n", cursor, p.prefix, p.description)
		}
		return tea.NewView(s)
	case messageScreen:
		s := "Write Message\n\n"
		s += m.choicedPrefix
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
		prefixDetails: []PrefixDetail{
			{"fix", "a commit that fixes a bug."},
			{"feat", "a commit that adds new functionality"},
			{"build", "changes that affect the build system or external dependencies."},
			{"chore", "other changes that don't modify src or test files."},
			{"ci", "changes to our CI configuration files and scripts."},
			{"docs", "a commit that adds or improves a documentation."},
			{"perf", "a commit that improves performance, without functional changes."},
			{"refactor", "a code change that neither fixes a bug nor adds a feature."},
			{"revert", "reverts a previous commit."},
			{"style", "changes that do not affect the meaning of the code."},
			{"test", "adding missing tests or correcting existing tests."},
		},
		choicedPrefix:        "",
		currentScreen: prefixScreen,
		textInput:     ti,
	}
}
