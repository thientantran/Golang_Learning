package main

import (
	"Food-delivery/hellopb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	hellopb.UnimplementedHelloServiceServer
}

func (*server) Hello(ctx context.Context, request *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	name := request.Name
	reponse := &hellopb.HelloResponse{
		Greeting: "Hi " + name,
	}
	return reponse, nil
}

func main() {
	address := "0.0.0.0:50051"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Err %v: ", err)
	}

	fmt.Printf("Server is listening on %v...:", address)

	s := grpc.NewServer()
	hellopb.RegisterHelloServiceServer(s, &server{})

	s.Serve(lis)

}
