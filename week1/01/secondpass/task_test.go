package secondpass

import (
	"errors"
	"testing"
)

func TestRewrite(t *testing.T) {
	item, err := NewTask("rewrite-1", "rewrite Task without notes")
	if err != nil {
		t.Fatalf("NewTask() returned an unexpected error: %v", err)
	}

	if err := item.Rename("check pointer receiver"); err != nil {
		t.Fatalf("Rename() returned an unexpected error: %v", err)
	}
	item.MarkDone()

	if item.Title != "check pointer receiver" || !item.IsDone() {
		t.Fatalf("unexpected task after changes: %+v", item)
	}
}

func TestRewriteRejectsEmptyRename(t *testing.T) {
	item, err := NewTask("rewrite-1", "original")
	if err != nil {
		t.Fatalf("NewTask() returned an unexpected error: %v", err)
	}

	if err := item.Rename(""); !errors.Is(err, ErrTitleRequired) {
		t.Fatalf("Rename() error = %v, want %v", err, ErrTitleRequired)
	}
	if item.Title != "original" {
		t.Fatalf("title changed on failed rename: %q", item.Title)
	}
}
