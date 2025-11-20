package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "hanifanmoha.com/grpc-1/hello"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello, " + req.GetName(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})

	log.Println("Server running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
