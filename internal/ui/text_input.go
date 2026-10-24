package ui

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

var ErrCancelled = errors.New("cancelled")

type TextInputModal struct {
	input     textinput.Model
	done      bool
	cancelled bool
	title     string
}

func NewInputModal(title, placeholder string) TextInputModal {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetWidth(len(ti.Placeholder) + 5)
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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m TextInputModal) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	return tea.NewView(fmt.Sprintf(
		"%s\n\n%s\n\n(press Enter to confirm)",
		m.title,
		m.input.View(),
	))
}

func (m TextInputModal) Run() (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	fm := finalModel.(TextInputModal)
	if fm.cancelled {
		return "", ErrCancelled
	}
	return fm.Value(), nil
}
