package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh"

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	status   int
	err      error
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		choices:  []string{"Buy carrots", "Buy Celery", "make soup"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit


	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": // WHO KNEW YOU COULD PUT MULTIPLE CHOICES IN A CASE STATEMENT
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	s += "\nPress q to quit.\n"

	return s
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

type (
	statusMsg int
	errMsg    struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }
