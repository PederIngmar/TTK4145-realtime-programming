// Use `go run foo.go` to run your program

package main

import (
	"fmt"
	"runtime"
)

func counter(inc_channel, dec_channel, done_channel, return_channel chan int) {
	i := 0
	doneflag := 0
	for {
		select {
		case <-inc_channel:
			i++
		case <-dec_channel:
			i--
		case <-done_channel:
			doneflag++
			if doneflag >= 2 {
				return_channel <- i
				return
			}
		}
	}
}

func incrementing(inc_channel, done_channel chan int) {
	for j := 0; j < 1000000; j++ {
		inc_channel <- 1
	}
	done_channel <- 1
}

func decrementing(dec_channel, done_channel chan int) {
	for j := 0; j < 1000000; j++ {
		dec_channel <- 1
	}
	done_channel <- 1
}

func main() {
	runtime.GOMAXPROCS(2)
	inc_channel := make(chan int)
	dec_channel := make(chan int)
	done_channel := make(chan int)
	return_channel := make(chan int)

	go incrementing(inc_channel, done_channel)
	go decrementing(dec_channel, done_channel)
	go counter(inc_channel, dec_channel, done_channel, return_channel)

	i := <-return_channel
	fmt.Println(i)
}
