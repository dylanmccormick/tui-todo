package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	name        string
	assignee    string
	description string
	status      status
	tags        []string
}

type tasksState struct {
	current  task
	cursor   int
	selected map[int]struct{}
	tasks    []task
}

type editTaskState struct {
	focusIndex int
	inputs     []textinput.Model
}

func (m model) TasksPageUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m.NewEditSwitch()
		case "D":
			// delete task
		case "up", "k":
			if m.state.tasks.cursor > 0 {
				m.state.tasks.cursor--
			}
		case "down", "j":
			if m.state.tasks.cursor < len(m.state.tasks.tasks)-1 {
				m.state.tasks.cursor++
			}
		case "enter", " ":
			_, ok := m.state.tasks.selected[m.state.tasks.cursor]
			if ok {
				delete(m.state.tasks.selected, m.state.tasks.cursor)
			} else {
				m.state.tasks.selected[m.state.tasks.cursor] = struct{}{}
			}

		}
	}

	return m, nil
}

func (m model) EditTasksUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	return m, cmd
}

func (m model) NewEditSwitch() (model, tea.Cmd) {
	m.state.tasks.current = task{}
	m = m.SwitchPage(editTaskPage)
	return m, nil
}

func (m model) TasksView() string {
	s := "My List of Tasks:\n\n"

	for i, task := range m.state.tasks.tasks {
		cursor := " "
		if m.state.tasks.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.state.tasks.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, task.name)
	}

	s += "\nPress n to add a new task."
	s += "\nPress q to quit.\n"

	return s
}

func (m model) EditTasksView() string {
	s := "Edit task here:\n"
	for i := range m.state.editTask.inputs {
		s += m.state.editTask.inputs[i].View()
		if i < len(m.state.editTask.inputs)-1 {
			s += "\n"
		}
	}
	s += "\nPress q to quit.\n"
	return s
}
