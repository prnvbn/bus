package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModal struct {
	input textinput.Model
	done  bool
	title string
}

func NewInputModal(title, placeholder string) TextInputModal {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = len(ti.Placeholder) + 5
	ti.Focus()

	return TextInputModal{input: ti, title: title}
}

func (m TextInputModal) Value() string {
	return m.input.Value()
}

func (m TextInputModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m TextInputModal) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n(press Enter to confirm)",
		m.title,
		m.input.View(),
	)
}

func (m TextInputModal) Run() (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	return finalModel.(TextInputModal).Value(), nil
}
