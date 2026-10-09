package main

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestMemoryNotifier(t *testing.T) {
	memory := NewMemoryNotifier()
	if err := SendWelcome(memory, "Alice"); err != nil {
		t.Fatal(err)
	}
	if err := SendNotification(memory, "练习已完成"); err != nil {
		t.Fatal(err)
	}
	want := []string{"欢迎你，Alice", "练习已完成"}
	if len(memory.Messages) != 2 || !reflect.DeepEqual(memory.Messages, want) {
		t.Fatalf("messages = %v, want %v", memory.Messages, want)
	}
}

func TestConsoleNotifier(t *testing.T) {
	var output bytes.Buffer
	if err := SendNotification(ConsoleNotifier{Writer: &output}, "hello"); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "控制台通知: hello\n" {
		t.Fatalf("console output = %q", got)
	}
}

func TestEmptyMessages(t *testing.T) {
	var memory MemoryNotifier
	var output bytes.Buffer
	for _, notifier := range []Notifier{&memory, ConsoleNotifier{Writer: &output}} {
		for _, message := range []string{"", " \t\n"} {
			if err := SendNotification(notifier, message); !errors.Is(err, ErrEmptyMessage) {
				t.Fatalf("business function error = %v", err)
			}
			if err := notifier.Notify(message); !errors.Is(err, ErrEmptyMessage) {
				t.Fatalf("direct Notify error = %v", err)
			}
		}
	}
	if len(memory.Messages) != 0 || output.Len() != 0 {
		t.Fatal("empty messages should not be saved or printed")
	}
}

type failingNotifier struct{ err error }

func (n failingNotifier) Notify(string) error { return n.err }

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func TestSendErrorsPropagate(t *testing.T) {
	want := errors.New("send failed")
	for _, notifier := range []Notifier{failingNotifier{want}, ConsoleNotifier{Writer: failingWriter{want}}} {
		if err := SendNotification(notifier, "hello"); !errors.Is(err, want) {
			t.Fatalf("send error = %v, want %v", err, want)
		}
		if err := SendWelcome(notifier, "Alice"); !errors.Is(err, want) {
			t.Fatalf("welcome error = %v, want %v", err, want)
		}
	}
}
