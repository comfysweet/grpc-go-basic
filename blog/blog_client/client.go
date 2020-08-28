package main

import (
	"context"
	"fmt"
	"github.com/comfysweet/grpc-go-basic/blog/blogpb"
	"google.golang.org/grpc"
	"io"
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

	// create blog
	c := blogpb.NewBlogServiceClient(conn)
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
	blogId := res.GetBlog().GetId()

	// read blog
	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogId,
	}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		log.Fatalf("Unexpected error: %v", readBlogErr)
	}
	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update blog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "KsZ2",
		Title:    "First blog (edited)",
		Content:  "Content of the blog with new articles",
	}
	updateBlogRes, updateBlogErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if updateBlogErr != nil {
		log.Fatalf("Unexpected error: %v", updateBlogErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateBlogRes)

	// delete blog
	deleteBlogRes, deleteBlogErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogId,
	})
	if deleteBlogErr != nil {
		log.Fatalf("Unexpected error: %v", deleteBlogErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteBlogRes)

	// list Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected error: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
