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
	Border(lipgloss.RoundedBorder(), true).
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
	List         list.Model
	CurrPath     string
	Width        int
	Height       int
	Hidden       bool
	ShortcutModal ShortcutModal 
}

func NewModel(l list.Model, path string) Model {
	shortcutItems := []list.Item{
		item{title: "Documents", desc: "📁", isDir: true, path: "/Users/lou/Documents"},
		item{title: "Desktop", desc: "📁", isDir: true, path: "/Users/lou/Desktop"},
		item{title: "EPL", desc: "📁", isDir: true, path: "/Users/lou/EPL"},
		item{title: "Cancel", desc: "", isDir: false, path: ""},
	}
	return Model{
		List:         l,
		CurrPath:     path,
		Hidden:       true,
		ShortcutModal: NewShortcutModal(shortcutItems, 30, 10),
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
	Back:  key.NewBinding(key.WithKeys( "left"), key.WithHelp("←", "go up directory")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
	Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "open directory")),
	Hide:  key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "toggle hidden files")),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.ShortcutModal.Active {
		var cmd tea.Cmd
		m.ShortcutModal, cmd = m.ShortcutModal.Update(msg)

		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.ShortcutModal.Active = false
			case "enter":
				if selected, ok := m.ShortcutModal.List.SelectedItem().(item); ok {
					if selected.title == "Cancel" {
						m.ShortcutModal.Active = false
					} else if selected.path != "" {
						m.CurrPath = selected.path
						m.ReloadDir()
						m.ShortcutModal.Active = false
					}
				}
			}
		}
		return m, cmd
	}

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

		case msg.String() == "s": // Show modal on 's' key
			m.ShortcutModal.Active = true
			return m, nil
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
    mainView := docStyle.Width(m.Width).Height(m.Height).Render(m.List.View())
    if m.ShortcutModal.Active {
        modal := m.ShortcutModal.View()
        centeredModal := lipgloss.Place(
            m.Width, m.Height,
            lipgloss.Center, lipgloss.Center,
            modal,
        )
        return mainView + "\n" + centeredModal
    }
    return mainView
}