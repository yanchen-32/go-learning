package main

import (
	"fmt"
	"log"

	"github.com/yanchen-32/go-learning/task"
)

func main() {
	item, err := task.New("task-1", "learn Go basics")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(item.Summary())
	if err := item.Rename("practice Go methods"); err != nil {
		log.Printf("rename task: %v", err)
		return
	}
	item.MarkDone()
	fmt.Println(item.Summary())

	// 调用方接收并处理 Rename 返回的校验错误。
	if err := item.Rename(""); err != nil {
		fmt.Printf("rename rejected: %v\n", err)
	}
}
