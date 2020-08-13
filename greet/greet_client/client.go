package main

import (
	"context"
	"fmt"
	"github.com/comfysweet/grpc-go-basic/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not to connect: %e", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do unary rpc")
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Ks",
		LastName:  "Z",
	}}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet grpc: %e", err)
	}
	log.Printf("Response from greet: %v", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do server streaming rpc")
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{
		FirstName: "Ks",
		LastName:  "Z",
	}}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet grpc: %e", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we have received all messages
			break
		}
		if err != nil {
			log.Fatalf("error while read server stream: %e", err)
		}
		log.Printf("received message: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do client streaming rpc\n")

	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Ks",
			LastName:  "Z",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "M",
			LastName:  "K",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "S",
			LastName:  "M",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "S",
			LastName:  "P",
		}},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		if err := stream.Send(req); err != nil {
			log.Fatalf("error while sending stream request: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from server: %v", err)
	}
	fmt.Printf("Received response: %v", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do BiDi streaming rpc\n")

	requests := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Ks",
			LastName:  "Z",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "M",
			LastName:  "K",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "S",
			LastName:  "M",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "S",
			LastName:  "P",
		}},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
	}

	waitc := make(chan struct{})
	//send a bunch of messages to the client
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending a message: %v\n", req)
			if err := stream.Send(req); err != nil {
				log.Fatalf("error while sending request to the client: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("error closing sending: %v", err)
		}
	}()
	//received bunch of messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while received bunch of responses")
				break
			}
			fmt.Printf("Received: %v\n", res)
		}
		close(waitc)
	}()

	//block until everyone is done
	<-waitc
}
