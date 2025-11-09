
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"fs/ui" 
)

func main() {
	startPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current dir:", err)
		os.Exit(1)
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			ui.Keys.Enter,
			ui.Keys.Right,
			ui.Keys.Back,
			ui.Keys.Quit,
			ui.Keys.Hide,
		}
	}

	delegate.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			{ui.Keys.Enter, ui.Keys.Right, ui.Keys.Back, ui.Keys.Hide},
			{ui.Keys.Quit},
		}
	}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Files"

	m := ui.NewModel(l, startPath)
	m.ReloadDir()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
