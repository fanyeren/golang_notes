package main

import (
	"github.com/ActiveState/tail"
	"fmt"
)


func main() {
	//t, err := tail.TailFile("/home/work/test.log", tail.Config{Follow: true})
	t, err := tail.TailFile("/home/work/test.log", tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		fmt.Println("error")
		return
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
