package main

//func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
//	log.Printf("client streaming started")
//	stream, err := client.SayHelloClientStreaming(context.Background())
//	if err != nil {
//		log.Fatalf("could not send names: %v", err)
//	}
//	for _, name := range names.Names {
//		req := &pb.HelloRequest{Name: name}
//		err := stream.Send(req)
//		if err != nil {
//			log.Fatalf("Error while sending: %v", err)
//		}
//		log.Printf("Send the request with name: %s", name)
//		time.Sleep(2 * time.Second)
//	}
//
//	recv, err := stream.CloseAndRecv()
//	log.Printf("client stream finished")
//	if err != nil {
//		log.Fatalf("Error while recving %v", err)
//	}
//	log.Printf("%v", recv.Messages)
//}
