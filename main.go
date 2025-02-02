package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	startTime   time.Time
	endTime     time.Time
	description string
}

type model struct {
	taskList  []task
	textinput textinput.Model
	err       error
}

type (
	errMsg error
)

func initialModel() model {
	textInput := textinput.New()
	textInput.Focus()
	textInput.Width = 25

	return model{
		taskList:  make([]task, 0),
		textinput: textInput,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			taskDescription := strings.TrimSpace(m.textinput.Value())
			currentTime := time.Now()
			if taskDescription == "**arrived" {
				m.taskList = append(m.taskList, task{
					description: taskDescription,
					startTime:   currentTime,
					endTime:     currentTime,
				})
			} else {
				m.taskList = append(m.taskList, task{
					description: taskDescription,
					startTime:   m.taskList[len(m.taskList)-1].endTime,
					endTime:     currentTime,
				})
			}

			m.textinput.SetValue("")
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func getHourFormat(taskTime time.Time) string {
	return taskTime.Format("15:04")
}

func (m model) View() string {
	var view string

	for _, task := range m.taskList {
		duration := task.endTime.Sub(task.startTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		formattedDuration := fmt.Sprintf("%02d h %02d min", hours, minutes)
		view += fmt.Sprintf("%s (%s-%s) %s\n", formattedDuration,
			getHourFormat(task.startTime),
			getHourFormat(task.endTime),
			task.description)
	}

	view += fmt.Sprintf(
		"Arrival Message%s\n\n%s",
		m.textinput.View(),
		"(esc to quit)",
	) + "\n"

	return view
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
