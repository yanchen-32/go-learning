// Package task 定义第一遍练习使用的 Task 类型及其行为。
package task

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyID    = errors.New("task ID cannot be empty")
	ErrEmptyTitle = errors.New("task title cannot be empty")
)

// Task 表示一个最小任务，包含编号、标题和完成状态。
type Task struct {
	ID    string
	Title string
	Done  bool
}

// New 校验必填字段并创建一个未完成的任务。
func New(id, title string) (*Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrEmptyID
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyTitle
	}

	return &Task{ID: id, Title: title}, nil
}

// Summary 只读取任务内容，因此使用 value receiver 即可。
func (t Task) Summary() string {
	return fmt.Sprintf("%s: %s (done=%t)", t.ID, t.Title, t.Done)
}

// Rename 会修改原任务的标题，因此使用 pointer receiver。
func (t *Task) Rename(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	t.Title = title
	return nil
}

// MarkDone 将原任务的状态修改为已完成。
func (t *Task) MarkDone() {
	t.Done = true
}
