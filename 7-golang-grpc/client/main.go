package main

import (
	"context"
	pb "golang-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const port = ":8000"

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("client streaming started")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("could not send names: %v", err)
	}
	for _, name := range names.Names {
		req := &pb.HelloRequest{Name: name}
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while sending: %v", err)
		}
		log.Printf("Send the request with name: %s", name)
		time.Sleep(2 * time.Second)
	}

	recv, err := stream.CloseAndRecv()
	log.Printf("client stream finished")
	if err != nil {
		log.Fatalf("Error while recving %v", err)
	}
	log.Printf("%v", recv.Messages)
}

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Akhil", "Alice", "Bob"},
	}

	//CallSayHello(client)

	client := pb.NewGreetServiceClient(conn)
	// callSayHelloServerStream(client, names)
	// callSayHelloClientStream(client, names)
	callBidirectionStream(client, names)
}
