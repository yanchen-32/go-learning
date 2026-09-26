package main

import (
	"fmt"
	"log"

	"github.com/yanchen-32/go-learning/secondpass"
)

func main() {
	item, err := secondpass.NewTask("rewrite-1", "rewrite Task from memory")
	if err != nil {
		log.Fatal(err)
	}

	if err := item.Rename("second pass finished"); err != nil {
		log.Fatal(err)
	}
	item.MarkDone()
	fmt.Printf("%+v, IsDone=%t\n", item, item.IsDone())
}
