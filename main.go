package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/myuron/ezgit/ui"
)

func main() {
	// entry point
	p := tea.NewProgram(ui.InitialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// crashes if an error occurs during a status update
	final := m.(ui.Model)
	if final.Err != nil {
		fmt.Fprintln(os.Stderr, final.Err)
		os.Exit(1)
	}
}
