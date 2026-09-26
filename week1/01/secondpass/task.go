package secondpass

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyID       = errors.New("task ID cannot be empty")
	ErrTitleRequired = errors.New("title Required cannot be empty")
)

type Task struct {
	ID    string
	Title string
	Done  bool
}

func NewTask(id, title string) (*Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrEmptyID
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}

	return &Task{ID: id, Title: title}, nil
}

func (t Task) Summary() string {
	return fmt.Sprintf("%s: %s (done=%t)", t.ID, t.Title, t.Done)
}

func (t *Task) Rename(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrTitleRequired
	}

	t.Title = title
	return nil
}
func (t Task) IsDone() bool {
	return t.Done
}
func (t *Task) MarkDone() {
	t.Done = true
}
