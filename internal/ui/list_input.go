package ui

import (
	"fmt"
	"iter"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type ListInputModel[T list.Item] struct {
	list      list.Model
	choice    T
	cancelled bool
}

func NewListInputModal[T list.Item](items iter.Seq[T]) ListInputModel[T] {
	li := make([]list.Item, 0)
	for item := range items {
		li = append(li, item)
	}

	l := list.New(li, list.NewDefaultDelegate(), 40, 0)

	return ListInputModel[T]{list: l}
}

func (m ListInputModel[T]) Init() tea.Cmd {
	return nil
}

func (m ListInputModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	preFilterState := m.list.FilterState()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
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

func (m ListInputModel[T]) View() tea.View {
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}

func (m ListInputModel[T]) Run() (T, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	sm := finalModel.(ListInputModel[T])
	if sm.cancelled {
		return sm.choice, ErrCancelled
	}

	return sm.choice, nil
}
