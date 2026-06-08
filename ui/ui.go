package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/myuron/ezgit/gitcontroller"
)

// Model ... A structure for managing the application's state
type Model struct {
	PrefixDetails []prefixDetail `json:"prefixDetails"`
	cursor        int
	choicedPrefix string
	currentScreen screen
	textInput     textinput.Model
	Err           error
}

// PrefixDetail ... A structure for conbining prefix and description
type prefixDetail struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
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
				if m.cursor < len(m.PrefixDetails)-1 {
					m.cursor++
				}
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "enter":
				m.choicedPrefix = m.PrefixDetails[m.cursor].Prefix
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
		for i, p := range m.PrefixDetails {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s: %s\n", cursor, p.Prefix, p.Description)
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

	m := Model{
		choicedPrefix: "",
		currentScreen: prefixScreen,
		textInput:     ti,
	}

	f, err := os.Open("./.ezgit/settings.json")
	if err != nil {
		f, err = os.Open("./settings.json")
		if err != nil {
			m.Err = err
			return m
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			m.Err = err
		}
	}()
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		m.Err = err
	}
	return m
}
