package main

import (
	"Food-delivery/hellopb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	fmt.Println("Hello, Client!")

	//opts := grpc.WithInsecure()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := hellopb.NewHelloServiceClient(cc)
	request := &hellopb.HelloRequest{Name: "Tan Dep Trai"}

	resp, _ := client.Hello(context.Background(), request)
	fmt.Printf("Receive Response => [%v]\n", resp.Greeting)
}
