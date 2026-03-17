package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DatabaseURL = "postgres://pguser123:pgpass123@localhost:5432/mydb?sslmode=disable"

func InitConnectionAndQuery(N int) {

	ctx := context.TODO()

	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		fmt.Printf("Unable to parse DATABASE_URL: %v\n", err)
		panic(err)
	}

	config.MaxConns = 100

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		panic(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			var timestamp time.Time
			err = conn.QueryRow(ctx, "SELECT NOW() FROM pg_sleep(3)").Scan(&timestamp)
			if err != nil {
				fmt.Printf("[%d] Error executing query: %v\n", i, err)
				return
			}
			fmt.Printf("[%d] %s\n", i, timestamp.String())
		}()
	}
	wg.Wait()
}

func main() {
	go func() {
		InitConnectionAndQuery(10)
	}()
	time.Sleep(60 * time.Second)
}
