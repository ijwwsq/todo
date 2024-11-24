package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type task struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"dateCreated"`
	IsDone      bool      `json:"isDone"`
	Deadline    time.Time `json:"deadline"`
}

var tasks = []task{
	{
		ID:          "1",
		Name:        "test todo",
		DateCreated: time.Now(),
		IsDone:      false,
		Deadline:    time.Time{},
	},
	{
		ID:          "2",
		Name:        "second todo",
		DateCreated: time.Now(),
		IsDone:      true,
		Deadline:    time.Time{},
	},
} // ? is needed?

func loadTasks() ([]task, error) {
	// yoink

	jsonFile, err := os.Open("tasks.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open tasks.json: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks.json: %w", err)
	}

	var tasks []task
	if err := json.Unmarshal(byteValue, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks.json: %w", err)
	}

	return tasks, nil
}

func saveTasks() error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile("tasks.json", data, 0644); err != nil {
		return fmt.Errorf("failed to save tasks.json: %w", err)
	}

	return nil
}

func showFullTasks() error {
	// yoink

	rows := [][]string{}

	for _, t := range tasks {
		dateCreated := t.DateCreated.Format("03:04 PM 01-02-2006")
		deadline := t.Deadline.Format("03:04 PM 01-02-2006")

		if deadline == "12:00 AM 01-01-0001" {
			deadline = "-"
		}

		var isDone string

		if t.IsDone {
			isDone = "yep"
		} else {
			isDone = "nope"
		}

		rows = append(rows, []string{t.ID, t.Name, dateCreated, isDone, deadline})
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle()
		}).
		Headers("ID", "TASK", "CREATED", "DONE", "DEADLINE").
		Rows(rows...)

	fmt.Println(t.Render())

	return nil
}

func showTasksShort() {
	var err error
	tasks, err = loadTasks()

	if err != nil {
		panic(err)
	}

	doneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	undoneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	taskStyle := lipgloss.NewStyle().Bold(true).MarginLeft(2)

	for _, t := range tasks {
		var checkmark string
		var taskDisplay string

		if t.IsDone {
			checkmark = "✓"
			taskDisplay = doneStyle.Render(t.Name)
		} else {
			checkmark = "•"
			taskDisplay = undoneStyle.Render(t.Name)
		}

		fmt.Printf("%s %s\n", checkmark, taskStyle.Render(taskDisplay))
	}
}

func updateTask(id string) {
	var err error
	tasks, err = loadTasks()

	if err != nil {
		panic(err)
	}

	for _, a := range tasks {
		if id == a.ID {
			fmt.Println("found, needs to be updated")
			return
		}
	}

	fmt.Println("not found")
}

func main() {
	saveTasks()
	var err error
	tasks, err = loadTasks()
	if err != nil {
		fmt.Printf("error loading tasks: %v\n", err)
	}

	showFullTasks()

	showTasksShort()

	updateTask("3")

	if err != nil { // !fix
		panic(err)
	}
}
