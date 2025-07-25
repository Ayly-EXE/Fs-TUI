package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"os/exec"


	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)

type item struct {
	title string
	desc  string
	isDir bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	currPath string
	width    int
	height   int
	hidden   bool
}

func (m *model) reloadDir() {
	entries, err := os.ReadDir(m.currPath)
	if err != nil {
		m.list.SetItems([]list.Item{item{title: "Error reading dir", desc: err.Error(), isDir: false}})
		return
	}

	var items []list.Item
	for _, entry := range entries {
	name := entry.Name()
	if m.hidden && strings.HasPrefix(name, ".") {
			continue 
	}

	desc := ""
	if entry.IsDir() {
		desc = "Directory"
		name += string(os.PathSeparator)
	} else {
		info, err := entry.Info()
		if err == nil {
			desc = fmt.Sprintf("File - %d bytes", info.Size())
		} else {
			desc = "File"
		}
	}
	items = append(items, item{title: name, desc: desc, isDir: entry.IsDir()})
}

	if len(items) == 0 {
		items = []list.Item{item{title: "(empty directory)", desc: "", isDir: false}}
	}

	m.list.SetItems(items)
	m.list.ResetSelected()
	m.list.Title = "Files in: " + m.currPath
}

func (m model) Init() tea.Cmd {
	return nil
}

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
		key.WithKeys("backspace", "left"),
		key.WithHelp("← / backspace", "go up directory"),
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Right):
			if selectedItem, ok := m.list.SelectedItem().(item); ok && selectedItem.isDir {
				dirName := strings.TrimRight(selectedItem.title, string(os.PathSeparator))
				newPath := filepath.Join(m.currPath, dirName)
				m.currPath = newPath
				m.reloadDir()
			}

		case key.Matches(msg, keys.Back):
			parent := filepath.Dir(m.currPath)
			if parent != m.currPath {
				m.currPath = parent
				m.reloadDir()
			}

		case key.Matches(msg, keys.Enter):
			selectedItem := m.list.SelectedItem().(item)
			path := filepath.Join(m.currPath, selectedItem.title)
			cmd := exec.Command("open", path) //Mac Os only : TODO -> Use Os type to open on Win and linux 
			if err := cmd.Start(); err != nil {
				fmt.Println("Failed to open:", err)
			}

		case key.Matches(msg, keys.Hide):
			m.hidden = !m.hidden
			m.reloadDir()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width - docStyle.GetHorizontalFrameSize()/2
		m.height = msg.Height - docStyle.GetVerticalFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Width(m.width).Height(m.height).Render(m.list.View())
}



func main() {
	startPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current dir:", err)
		os.Exit(1)
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	// Set keybindings help for the delegate
	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			keys.Enter,
			keys.Right,
			keys.Back,
			keys.Quit,
			keys.Hide,
		}
	}

	delegate.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			{keys.Enter, keys.Right, keys.Back, keys.Hide},
			{keys.Quit},
		}
	}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Files"

	m := model{
		list:     l,
		currPath: startPath,
		hidden: true,
	}

	m.reloadDir()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
