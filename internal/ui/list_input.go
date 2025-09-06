package ui

import (
	"iter"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListInputModel[T list.Item] struct {
	list   list.Model
	choice T
}

func NewListInputModal[T list.Item](items iter.Seq[T]) ListInputModel[T] {
	li := make([]list.Item, 0)
	for item := range items {
		li = append(li, item)
	}

	l := list.New(li, list.NewDefaultDelegate(), 40, 0)

	return ListInputModel[T]{list: l}
}

func (m ListInputModel[T]) Choice() T {
	return m.choice
}

func (m ListInputModel[T]) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m ListInputModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	preFilterState := m.list.FilterState()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if preFilterState != list.Filtering {
				if item, ok := m.list.SelectedItem().(T); ok {
					m.choice = item
				}
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListInputModel[T]) View() string {
	return m.list.View()
}
