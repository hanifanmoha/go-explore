package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "hanifanmoha.com/grpc-1/hello"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Did not connect:", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Hanifan"})
	if err != nil {
		log.Fatal("Error:", err)
	}

	log.Println("Response:", resp.GetMessage())
}
