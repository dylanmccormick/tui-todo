package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	tasksPage page = iota
	editTaskPage
)

var style = lipgloss.NewStyle().
	Width(100)

type model struct {
	page        page
	taskList    taskList
	taskMenu    taskMenu
	state       state
	pageChanged bool
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
	m.pageChanged = true
	m.page = page
	return m
}

func (m model) View() string {
	switch m.page {
	case tasksPage:
		return style.Render(m.TasksView())
	case editTaskPage:
		return style.Render(m.EditTasksView())
	}

	return ""
}

func baseModel() model {
	taskState := initTaskState()
	editTaskState := initEditTaskState()
	return model{state: state{
		tasks:    taskState,
		editTask: editTaskState,
	}}
}

func initTaskState() tasksState {
	tasks, err := ReadFromFile()
	if err != nil {
		tasks = []task{
			{
				Name:        "take out garbage",
				Description: "put the garbage in the government garbage can",
				Assignee:    "Dylan",
				Status:      NOT_STARTED,
				Tags:        []string{"chores"},
			},
			{
				Name:        "do the dishes",
				Description: "clean all of the dishes in the sink",
				Assignee:    "Dylan",
				Status:      DONE,
				Tags:        []string{"chores"},
			},
		}
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
		case "ctrl+c":
			m.SaveToFile()
			return m, tea.Quit
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
