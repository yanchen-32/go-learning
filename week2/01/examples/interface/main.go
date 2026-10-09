package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrEmptyMessage = errors.New("message cannot be empty")

// Notifier 只描述调用方真正需要的行为。
type Notifier interface {
	Notify(message string) error
}

// Writer 未设置时输出到标准输出；测试可以注入 buffer。
type ConsoleNotifier struct {
	Writer io.Writer
}

func (n ConsoleNotifier) Notify(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	writer := n.Writer
	if writer == nil {
		writer = os.Stdout
	}
	_, err := fmt.Fprintln(writer, "控制台通知:", message)
	return err
}

type MemoryNotifier struct {
	Messages []string
}

// 返回具体类型，既能作为 Notifier 使用，也能读取 Messages。
func NewMemoryNotifier() *MemoryNotifier {
	return &MemoryNotifier{Messages: make([]string, 0)}
}

// Notify 要保存消息，因此使用 pointer receiver。
func (n *MemoryNotifier) Notify(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	n.Messages = append(n.Messages, message)
	return nil
}

// SendNotification 依赖行为，并将发送错误返回给调用方。
func SendNotification(notifier Notifier, message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	return notifier.Notify(message)
}

// SendWelcome 复用依赖 interface 的业务函数。
func SendWelcome(notifier Notifier, name string) error {
	return SendNotification(notifier, "欢迎你，"+name)
}

// 编译期检查：如果类型没有完整实现 Notifier，这里会编译失败。
var _ Notifier = ConsoleNotifier{}
var _ Notifier = (*MemoryNotifier)(nil)

func main() {
	console := ConsoleNotifier{}
	if err := SendWelcome(console, "Alice"); err != nil {
		fmt.Println("发送失败:", err)
	}

	memory := NewMemoryNotifier()
	if err := SendWelcome(memory, "Bob"); err != nil {
		fmt.Println("发送失败:", err)
	}
	if err := SendNotification(memory, "练习已完成"); err != nil {
		fmt.Println("发送失败:", err)
	}
	if err := SendNotification(memory, "  "); err != nil {
		fmt.Println("空消息发送失败:", err)
	}
	fmt.Println("成功发送次数:", len(memory.Messages))
	fmt.Println("内存中的消息:", memory.Messages)
}
