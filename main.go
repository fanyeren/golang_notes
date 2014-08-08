package main

import (
	"fmt"
	"time"
	"runtime"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil"
)

func test(c chan bool) {
	c<-true
}

func send (c chan<- int) {
}

func recv (c <-chan int) {
}

var sem = make (chan int, 2)
var over = make (chan bool)

func worker(i int) {
	sem <- 1
	fmt.Println(time.Now())
	<-sem

	if i == 9 { over <- true }
}

func main() {
	runtime.GOMAXPROCS(2)
	c := make (chan int)
	var procs, err = ps.Processes()
	if err == nil {
		for _, value := range procs {
			var proc = value.Pid()
			var pproc = value.PPid()
			var procname = value.Executable()
			fmt.Println(proc, pproc, procname)
		}
	}

        v, _ := gopsutil.VirtualMemory()

        // almost every return value is struct
        fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

        // convert to JSON. String() is also implemented
        fmt.Println(v)

	hv, _ := gopsutil.HostInfo()
        fmt.Println(hv.VirtualizationSystem)

	go func() {
		for i := 0; i < 20 ; i++ {
			c<-i
		}
		close(c)
	}()

	//for v := range c {
	//	println(v)
	//}

	for {
		if v, ok := <-c; ok {
			println(v)
		} else {
			break;
		}
	}


	for i:= 0; i < 10; i++ {
		go worker(i)
	}

	<-over
}
