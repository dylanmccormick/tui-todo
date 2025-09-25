package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	TasksView AppState = iota
	AddTaskView
)

type model struct {
	currentView AppState
	taskList    taskList
	taskMenu    taskMenu
}

type taskMenu struct {
	focusIndex int
	inputs     []textinput.Model
}

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
	p := tea.NewProgram(baseModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) View() string {
	switch m.currentView {
	case TasksView:
		return m.taskList.View()
	case AddTaskView:
		return m.taskMenu.View()
	}

	return ""
}

func baseMenu() taskMenu {
	tm := taskMenu{
		inputs: make([]textinput.Model, 1),
	}
	var t textinput.Model
	for i := range tm.inputs {
		t = textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()
		}
		tm.inputs[i] = t
	}
	return tm
}

func baseModel() model {
	baseList := baseList()
	taskMenu := baseMenu()
	return model{taskList: baseList, taskMenu: taskMenu}
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
		case "ctrl+c", "q": // WHO KNEW YOU COULD PUT MULTIPLE CHOICES IN A CASE STATEMENT
			return m, tea.Quit
		case "t":
			m.currentView = TasksView
			return m, nil
		case "n":
			m.currentView = AddTaskView
			return m, nil
		}
	}

	var cmd tea.Cmd
	var ok bool
	switch m.currentView {
	case TasksView:
		var tl any
		tl, cmd = m.taskList.Update(msg)
		m.taskList, ok = tl.(taskList)
		if !ok {
			panic("unexpected type")
		}
	case AddTaskView:
		var tv any
		tv, cmd = m.taskMenu.Update(msg)
		m.taskMenu, ok = tv.(taskMenu)
		if !ok {
			panic("unexpected type")
		}
	}
	return m, cmd
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
		case "n":
			tl.createNewTask()

		}
	}
	return tl, nil
}

func (tl taskList) createNewTask() {
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

	s += "\nPress n to add a new task."
	s += "\nPress q to quit.\n"

	return s
}

func (tm taskMenu) Init() tea.Cmd {
	return textinput.Blink
}

func (tm taskMenu) View() string {
	s := "Edit task here:\n"
	for i := range tm.inputs {
		s += tm.inputs[i].View()
		if i < len(tm.inputs)-1 {
			s += "\n"
		}
	}
	s += "\nPress q to quit.\n"
	return s
}

func (tm taskMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type){
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			cmds := make([]tea.Cmd, len(tm.inputs))
			for i:=0; i<=len(tm.inputs)-1; i++{

			}
			return tm, tea.Batch(cmds...)
		}
	}
	cmd := tm.updateInputs(msg)
	return tm, cmd
}

func (tm *taskMenu) updateInputs(msg tea.Msg) (tea.Cmd){
	cmds := make([]tea.Cmd, len(tm.inputs))

	for i := range tm.inputs {
		tm.inputs[i], cmds[i] = tm.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
