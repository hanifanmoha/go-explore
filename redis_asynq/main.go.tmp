package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

const TaskFib = "task_fib"

type FibPayload struct {
	N int
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func NewFibTask(n int) (*asynq.Task, error) {
	payload, err := json.Marshal(FibPayload{N: n})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskFib, payload), nil
}

func HandleFibonacciTask(ctx context.Context, t *asynq.Task) error {

	var p FibPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	fmt.Printf("Running task for n=%d\n", p.N)
	result := fib(p.N)
	fmt.Printf("Fibonacci of %d is %d\n", p.N, result)

	err := os.WriteFile(fmt.Sprintf("./out/%d.txt", p.N), []byte(fmt.Sprintf("%d", result)), 0644)
	if err != nil {
		return err
	}

	return nil
}

func enqueue(client *asynq.Client, n int) error {
	task, err := NewFibTask(n)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)
	if err != nil {
		fmt.Println("Could not enqueue task:", err)
		return err
	}

	fmt.Printf("Enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}

func main() {

	redisClient := asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	}

	client := asynq.NewClient(redisClient)
	defer client.Close()

	enqueue(client, 11)
	enqueue(client, 51)
	enqueue(client, 21)
	enqueue(client, 44)
	enqueue(client, 31)

	runServer(&redisClient)

}

func runServer(redisClient *asynq.RedisClientOpt) {

	srv := asynq.NewServer(redisClient, asynq.Config{
		Concurrency: 2,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskFib, HandleFibonacciTask)

	go func() {
		if err := srv.Start(mux); err != nil {
			fmt.Println("Could not start server:", err)
		}
	}()

	time.Sleep(60 * time.Second)

	srv.Shutdown()
	fmt.Println("Program finished!")
}
