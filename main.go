package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/myuron/ezgit/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if final, ok := m.(ui.Model); ok && final.Err != nil {
		fmt.Fprintln(os.Stderr, final.Err)
		os.Exit(1)
	}
}
