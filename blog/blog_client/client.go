package main

import (
	"context"
	"fmt"
	"github.com/comfysweet/grpc-go-basic/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not to connect: %e", err)
	}
	defer conn.Close()

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "KsZ",
		Title:    "First blog",
		Content:  "Content of the blog",
	}

	c := blogpb.NewBlogServiceClient(conn)
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
}
