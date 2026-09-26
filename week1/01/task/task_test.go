package task

import (
	"errors"
	"testing"
)

func TestTaskLifecycle(t *testing.T) {
	item, err := New("task-1", "learn struct")
	if err != nil {
		t.Fatalf("New() returned an unexpected error: %v", err)
	}

	if err := item.Rename("learn methods"); err != nil {
		t.Fatalf("Rename() returned an unexpected error: %v", err)
	}
	item.MarkDone()

	if item.Title != "learn methods" || !item.Done {
		t.Fatalf("unexpected task after changes: %+v", item)
	}
}

func TestRenameRejectsBlankTitle(t *testing.T) {
	item, err := New("task-1", "keep this title")
	if err != nil {
		t.Fatalf("New() returned an unexpected error: %v", err)
	}

	err = item.Rename("   ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("Rename() error = %v, want %v", err, ErrEmptyTitle)
	}
	if item.Title != "keep this title" {
		t.Fatalf("Rename() changed the title after an error: %q", item.Title)
	}
}

func TestNewValidatesRequiredFields(t *testing.T) {
	if _, err := New("", "title"); !errors.Is(err, ErrEmptyID) {
		t.Fatalf("New() error = %v, want %v", err, ErrEmptyID)
	}
	if _, err := New("task-1", ""); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("New() error = %v, want %v", err, ErrEmptyTitle)
	}
}
