## ✅ 1. Install Required Tools
**Install protoc**

Mac:
```
brew install protobuf
```
Install Go plugins
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Add to PATH (if needed):
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

## ✅ 2. Create a .proto File

Example: `proto/hello.proto`
```
syntax = "proto3";

package hello;
option go_package = "example.com/myapp/hello";

service HelloService {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

## ✅ 3. Generate Go Code

Run:
```
protoc \
  --go_out=. \
  --go-grpc_out=. \
  proto/hello.proto
```

This generates two files inside /hello:

- `hello.pb.go`
- `hello_grpc.pb.go`

## ✅ 4. Implement the gRPC Server

`server/main.go`:

```
package main

import (
	"context"
	"log"
	"net"

	pb "example.com/myapp/hello"
	"google.golang.org/grpc"
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
```

Run:
```
go run server/main.go
```

## ✅ 5. Implement the gRPC Client

`client/main.go`:
```
package main

import (
	"context"
	"log"
	"time"

	pb "example.com/myapp/hello"
	"google.golang.org/grpc"
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
```

Run:
```
go run client/main.go
```

## 🎉 Result

Server prints:
```
Server running on port 50051
```

Client prints:
```
Response: Hello, Hanifan
```