package main

import (
	"os"
	"fmt"
)

func main () {
	args := os.Args
	arg_num := len(args)

	fmt.Printf("args_num is %d\n", arg_num)

	for i := 0; i < arg_num; i++ {
		fmt.Println(args[i])
	}

	environ := os.Environ()
	for i := range environ {
		fmt.Println(environ[i])
	}

	logname := os.Getenv("LOGNAME")
	fmt.Printf("logname is %s\n", logname)
}
