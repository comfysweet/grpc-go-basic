package main

import (
	"context"
	"fmt"
	"github.com/comfysweet/grpc-go-basic/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with: %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client canceled the request
			fmt.Println("The client canceled the request")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{Result: result}
	return res, nil
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Print("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while receiving client request: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		res := "Hello " + firstName

		if errSend := stream.Send(&greetpb.GreetEveryoneResponse{Result: res,}); errSend != nil {
			log.Fatalf("error sending response to client")
			return errSend
		}
	}
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Print("Long greet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading a stream request
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		result += "Hello " + req.GetGreeting().GetFirstName() + "\n"
	}

}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet many time function was invoked with: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{Result: result}
		if err := stream.Send(res); err != nil {
			log.Fatalf("Failed to serve: %v", err)
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
