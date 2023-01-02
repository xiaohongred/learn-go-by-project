package main

import (
	pb "golang-grpc/proto"
	"log"
	"time"
)

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList,
	steram pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("got request with name: %v", req.Names)
	for _, name := range req.Names {
		res := &pb.HelloResponse{Message: "Hello" + name}
		if err := steram.Send(res); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
