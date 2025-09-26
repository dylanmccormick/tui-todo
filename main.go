package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	tasksPage page = iota
	editTaskPage
)

type model struct {
	page     page
	taskList taskList
	taskMenu taskMenu
	state    state
}

type state struct {
	tasks    tasksState
	editTask editTaskState
}

type taskMenu struct {
	focusIndex int
	inputs     []textinput.Model
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
	p := tea.NewProgram(baseModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) SwitchPage(page page) model {
	m.page = page
	return m
}

func (m model) View() string {
	switch m.page {
	case tasksPage:
		return m.TasksView()
	case editTaskPage:
		return m.EditTasksView()
	}

	return ""
}

func baseModel() model {
	taskState := initTaskState()
	return model{state: state{
		tasks: taskState,
	},}
}

func initTaskState() tasksState {
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
	return tasksState{tasks: tasks, selected: make(map[int]struct{})}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: We can't enter in the form when we're always scanning for "t" and "n". I can't type any word with those letters while using the form
		// maybe we should just use ctrl+n, ctrl+t??
		// or we have a state that is something like "inputting text" and the only command is esc/ tab (esc to go back, tab/shift+tab to navigate
		case "ctrl+c": // WHO KNEW YOU COULD PUT MULTIPLE CHOICES IN A CASE STATEMENT
			return m, tea.Quit
		case "esc":
			// TODO: This will need to go back on other pages. But for now do nothing
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.page {
	case tasksPage:
		m, cmd = m.TasksPageUpdate(msg)
	case editTaskPage:
		m, cmd = m.EditTasksUpdate(msg)
	}
	return m, cmd
}

