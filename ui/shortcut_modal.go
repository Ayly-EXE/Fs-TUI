package ui

import (
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var shortcutModalStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder(), true).
    BorderForeground(lipgloss.Color("205")).
    Padding(1, 2)

type ShortcutModal struct {
    List   list.Model
    Active bool
    Width  int
    Height int
}

func NewShortcutModal(items []list.Item, width, height int) ShortcutModal {
    l := list.New(items, list.NewDefaultDelegate(), width, height)
    l.Title = "Shortcuts"
    return ShortcutModal{
        List:   l,
        Active: false,
        Width:  width,
        Height: height,
    }
}

func (m ShortcutModal) Init() tea.Cmd {
    return nil
}

func (m ShortcutModal) Update(msg tea.Msg) (ShortcutModal, tea.Cmd) {
    if !m.Active {
        return m, nil
    }
    var cmd tea.Cmd
    m.List, cmd = m.List.Update(msg)
    return m, cmd
}

func (m ShortcutModal) View() string {
    if !m.Active {
        return ""
    }
    return shortcutModalStyle.Width(m.Width).Height(m.Height).Render(m.List.View())
}