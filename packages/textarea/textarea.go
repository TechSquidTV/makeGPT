package textarea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func New(placeholder string) (string, error) {
	p := tea.NewProgram(initialModel(placeholder))

	m, err := p.Run()
	if err != nil {
		return "", err
	}

	model, ok := m.(textAreaModel) // Corrected type assertion
	if !ok {
		return "", fmt.Errorf("could not assert model type")
	}

	return model.textarea.Value(), nil
}

type textAreaModel struct {
	textarea textarea.Model
	err      error
}

func initialModel(placeholder string) textAreaModel {
	ti := textarea.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 0 // Set to 0 for no limit

	return textAreaModel{
		textarea: ti,
		err:      nil,
	}
}

func (m textAreaModel) Init() tea.Cmd {
	return nil
}

func (m textAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlS:
			saveContent(m.textarea.Value())
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m textAreaModel) View() string {
	return fmt.Sprintf(
		"Write your text here. (ctrl+s to save and exit, ctrl+c to quit without saving)\n\n%s",
		m.textarea.View(),
	) + "\n\n"
}

func saveContent(content string) {
	log.Info("Saved content")
}
