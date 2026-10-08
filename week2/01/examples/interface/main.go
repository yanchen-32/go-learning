package main

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyMessage = errors.New("message cannot be empty")

// Notifier 只描述调用方真正需要的行为。
type Notifier interface {
	Notify(message string) error
}

type ConsoleNotifier struct{}

func (ConsoleNotifier) Notify(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	fmt.Println("控制台通知:", message)
	return nil
}

type MemoryNotifier struct {
	Messages []string
}

// Notify 要保存消息，因此使用 pointer receiver。
func (n *MemoryNotifier) Notify(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	n.Messages = append(n.Messages, message)
	return nil
}

// SendWelcome 依赖行为而不是具体通知器。
func SendWelcome(notifier Notifier, name string) error {
	return notifier.Notify("欢迎你，" + name)
}

// 编译期检查：如果类型没有完整实现 Notifier，这里会编译失败。
var _ Notifier = ConsoleNotifier{}
var _ Notifier = (*MemoryNotifier)(nil)

func main() {
	console := ConsoleNotifier{}
	if err := SendWelcome(console, "Alice"); err != nil {
		fmt.Println("发送失败:", err)
	}

	memory := &MemoryNotifier{}
	if err := SendWelcome(memory, "Bob"); err != nil {
		fmt.Println("发送失败:", err)
	}
	fmt.Println("内存中的消息:", memory.Messages)
}
