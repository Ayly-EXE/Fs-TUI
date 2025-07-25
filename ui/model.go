// File: ui/model.go
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

type Model struct {
	List     list.Model
	CurrPath string
	Width    int
	Height   int
	Hidden   bool
}

func NewModel(l list.Model, path string) Model {
	return Model{
		List:     l,
		CurrPath: path,
		Hidden:   true,
	}
}

func (m *Model) ReloadDir() {
	entries, err := os.ReadDir(m.CurrPath)
	if err != nil {
		m.List.SetItems([]list.Item{item{title: "Error reading dir", desc: err.Error(), isDir: false}})
		return
	}

	var items []list.Item
	for _, entry := range entries {
		name := entry.Name()
		if m.Hidden && strings.HasPrefix(name, ".") {
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

	m.List.SetItems(items)
	m.List.ResetSelected()
	m.List.Title = "Files in: " + m.CurrPath
}

var Keys = struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Right key.Binding
	Hide  key.Binding
}{
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Access file or dir")),
	Back:  key.NewBinding(key.WithKeys("backspace", "left"), key.WithHelp("← / backspace", "go up directory")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
	Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "open directory")),
	Hide:  key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "toggle hidden files")),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, Keys.Right):
			if selectedItem, ok := m.List.SelectedItem().(item); ok && selectedItem.isDir {
				dirName := strings.TrimRight(selectedItem.title, string(os.PathSeparator))
				newPath := filepath.Join(m.CurrPath, dirName)
				m.CurrPath = newPath
				m.ReloadDir()
			}

		case key.Matches(msg, Keys.Back):
			parent := filepath.Dir(m.CurrPath)
			if parent != m.CurrPath {
				m.CurrPath = parent
				m.ReloadDir()
			}

		case key.Matches(msg, Keys.Enter):
			selectedItem := m.List.SelectedItem().(item)
			path := filepath.Join(m.CurrPath, selectedItem.title)
			cmd := exec.Command("open", path) // macOS only — use runtime.GOOS to generalize
			if err := cmd.Start(); err != nil {
				fmt.Println("Failed to open:", err)
			}

		case key.Matches(msg, Keys.Hide):
			m.Hidden = !m.Hidden
			m.ReloadDir()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.Width = msg.Width - docStyle.GetHorizontalFrameSize()/2
		m.Height = msg.Height - docStyle.GetVerticalFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Width(m.Width).Height(m.Height).Render(m.List.View())
}
