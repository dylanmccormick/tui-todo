package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh"

type task struct {
	name        string
	assignee    string
	description string
	status      status
	tags        []string
}

type taskList struct {
	tasks    []task
	cursor   int
	selected map[int]struct{}
}

type status int

const (
	NOT_STARTED status = iota
	DOING
	DONE
)

func main() {
	p := tea.NewProgram(baseList())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func baseList() taskList {
	tasks := []task{
		{
			name:        "take out garbage",
			description: "put the garbage in the government garbage can",
			assignee:    "Dylan",
			status:      NOT_STARTED,
			tags:        []string{"chores"},
		},
		{
			name:        "do the dishes",
			description: "clean all of the dishes in the sink",
			assignee:    "Dylan",
			status:      DONE,
			tags:        []string{"chores"},
		},
	}
	return taskList{tasks: tasks, selected: make(map[int]struct{})}
}

func (tl taskList) Init() tea.Cmd {
	return nil
}

func (tl taskList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": // WHO KNEW YOU COULD PUT MULTIPLE CHOICES IN A CASE STATEMENT
			return tl, tea.Quit
		case "up", "k":
			if tl.cursor > 0 {
				tl.cursor--
			}
		case "down", "j":
			if tl.cursor < len(tl.tasks)-1 {
				tl.cursor++
			}
		case "enter", " ":
			_, ok := tl.selected[tl.cursor]
			if ok {
				delete(tl.selected, tl.cursor)
			} else {
				tl.selected[tl.cursor] = struct{}{}
			}

		}
	}
	return tl, nil
}

func (tl taskList) View() string {
	s := "My List of Tasks:\n\n"

	for i, task := range tl.tasks {
		cursor := " "
		if tl.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := tl.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, task.name)
	}

	s += "\nPress q to quit.\n"

	return s
}
