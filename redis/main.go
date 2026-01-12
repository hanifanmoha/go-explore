package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func val2Key(val int) string {
	return fmt.Sprintf("fib:%d", val)
}

func fibCache(ctx context.Context, rdb *redis.Client, n int) int {
	if n <= 1 {
		return n
	}

	val, err := rdb.Get(ctx, val2Key(n)).Result()
	if err == nil {
		cachedVal, _ := strconv.Atoi(val)
		return cachedVal
	}

	newVal := fibCache(ctx, rdb, n-1) + fibCache(ctx, rdb, n-2)
	rdb.Set(ctx, val2Key(n), strconv.Itoa(newVal), time.Second*60)
	return newVal
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       3,
	})

	defer rdb.Close()

	ctx := context.Background()

	n := 50
	// fib10 := fib(n)
	// fmt.Printf("Fibonacci of %d is (uncached): %d\n", n, fib10)
	fib10Cache := fibCache(ctx, rdb, n)
	fmt.Printf("Fibonacci of %d is (cached): %d\n", n, fib10Cache)

}
