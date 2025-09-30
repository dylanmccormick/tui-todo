package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	Name        string   `json:"name"`
	Assignee    string   `json:"assignee"`
	Description string   `json:"description"`
	Status      status   `json:"status"`
	Tags        []string `json:"tags"`
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

const (
	name = iota
	description
	assignee
	stat
	tag
)

func (m model) TasksPageUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m.NewEditSwitch()
		case "E":
			return m.EditSwitch()
		case "D":
			s := m.state.tasks.tasks
			m.state.tasks.tasks = append(s[:m.state.tasks.cursor], s[m.state.tasks.cursor+1:]...)
			if len(m.state.tasks.tasks) == 0 {
				m.state.tasks.cursor = 0
			} else {
				m.state.tasks.cursor = 1
			}
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
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.state.editTask.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			m.state.tasks.current = task{}
			m = m.SwitchPage(tasksPage)
			return m, nil
		case tea.KeyEnter:
			m.nextInput()
		case tea.KeyShiftTab:
			m.prevInput()
		case tea.KeyTab:
			m.nextInput()
		case tea.KeyCtrlS:
			m.SaveEditTask()
			m.state.tasks.current = task{}
			m = m.SwitchPage(tasksPage)
		}
		for i := range m.state.editTask.inputs {
			m.state.editTask.inputs[i].Blur()
		}
		m.state.editTask.inputs[m.state.editTask.focusIndex].Focus()
	}
	for i := range m.state.editTask.inputs {
		m.state.editTask.inputs[i], cmds[i] = m.state.editTask.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) NewEditSwitch() (model, tea.Cmd) {
	m.state.tasks.current = task{}
	m.state.editTask.inputs[name].SetValue(m.state.tasks.current.Name)
	m.state.editTask.inputs[description].SetValue(m.state.tasks.current.Description)
	m.state.editTask.inputs[assignee].SetValue(m.state.tasks.current.Assignee)
	m = m.SwitchPage(editTaskPage)
	return m, nil
}

func (m model) EditSwitch() (model, tea.Cmd) {
	m.state.tasks.current = m.state.tasks.tasks[m.state.tasks.cursor]
	m.state.editTask.inputs[name].SetValue(m.state.tasks.current.Name)
	m.state.editTask.inputs[description].SetValue(m.state.tasks.current.Description)
	m.state.editTask.inputs[assignee].SetValue(m.state.tasks.current.Assignee)
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

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, task.Name)
	}

	s += "\nPress n to add a new task."
	s += "\nPress d to delete a task."

	return s
}

func (m model) EditTasksView() string {
	s := "Edit task here:\n"
	s += fmt.Sprintf(
		`	Task Name: %s
	Description: %s
	Assignee: %s
		`, m.state.editTask.inputs[name].View(),
		m.state.editTask.inputs[description].View(),
		m.state.editTask.inputs[assignee].View()) + "\n"
	s += "Press ctrl+s to save and exit\n"
	s += "\nPress esc to exit without saving\n"
	return s
}

func (m *model) nextInput() {
	m.state.editTask.focusIndex = (m.state.editTask.focusIndex + 1) % len(m.state.editTask.inputs)
}

func (m *model) prevInput() {
	m.state.editTask.focusIndex--
	// Wrap around
	if m.state.editTask.focusIndex < 0 {
		m.state.editTask.focusIndex = len(m.state.editTask.inputs) - 1
	}
}

func initEditTaskState() editTaskState {
	inputs := make([]textinput.Model, 3)
	inputs[name] = textinput.New()
	inputs[name].Focus()
	inputs[name].Prompt = ""
	inputs[name].Placeholder = "name"

	inputs[description] = textinput.New()
	inputs[description].Prompt = ""
	inputs[description].Placeholder = "description"

	inputs[assignee] = textinput.New()
	inputs[assignee].Prompt = ""
	inputs[assignee].Placeholder = "description"

	return editTaskState{
		inputs: inputs,
	}
}

func (m *model) SaveEditTask() {
	newTask := task{
		Name:        m.state.editTask.inputs[name].Value(),
		Description: m.state.editTask.inputs[description].Value(),
		Assignee:    m.state.editTask.inputs[assignee].Value(),
	}

	if m.state.tasks.current.Name == "" {
		m.state.tasks.tasks = append(m.state.tasks.tasks, newTask)
	} else {
		m.state.tasks.tasks[m.state.tasks.cursor] = newTask
	}
}
