package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(q chan int) {
	for {
		for i := 0; i < 5; i++ {
			value := rand.Intn(1000)
			fmt.Println("Producing...", value)
			q <- value
		}
		sleepTime := rand.Intn(3) + 2
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func consumer(q chan int) {
	for {
		time.Sleep(2 * time.Second)
		fmt.Println("Waiting to consume...")
		value := <-q
		fmt.Println("Consuming...", value)
	}
}

func main() {
	q := make(chan int, 10)

	go producer(q)
	go consumer(q)

	for {
		time.Sleep(1 * time.Second)
	}
}
