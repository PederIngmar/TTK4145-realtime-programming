package main

import (
	"fmt"
	"time"
)

func producer(buffered_channel chan int) {

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[producer]: pushing %d\n", i)
		buffered_channel <- i
	}
	//close(buffered_channel)
}

func consumer(buffered_channel chan int) {

	time.Sleep(1 * time.Second)
	for {
		i := <-buffered_channel
		fmt.Printf("[consumer]: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}

}

func main() {
	buffered_channel := make(chan int, 5)

	go consumer(buffered_channel)
	go producer(buffered_channel)

	select {}
}
