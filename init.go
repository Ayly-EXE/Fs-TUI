// File: init.go
package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Right key.Binding
	Hide  key.Binding
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Access file or dir"),
	),
	Back: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go up directory"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "open directory"),
	),
	Hide: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "toggle hidden files"),
	),
}
