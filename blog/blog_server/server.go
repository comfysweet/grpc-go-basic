package main

import (
	"context"
	"fmt"
	"github.com/comfysweet/grpc-go-basic/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var collection *mongo.Collection

type server struct {
}

func (s *server) ListBlog(request *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &blogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data: %v", err))
		}
		stream.Send(&blogpb.ListBlogResponse{
			Blog: dataToBlogPb(data),
		})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}
	return nil
}

func (s *server) DeleteBlog(ctx context.Context, request *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse id: %v", err))
	}
	filter := bson.D{{"_id", oid}}
	res, deleteErr := collection.DeleteOne(context.Background(), filter)
	if deleteErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot delete object in MongoDB: %v", deleteErr))
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog in MongoDB: %v", deleteErr))
	}
	return &blogpb.DeleteBlogResponse{
		BlogId: blogId,
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, request *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := request.GetBlog()
	blogId := blog.GetId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse id: %v", err))
	}
	data := &blogItem{}
	filter := bson.D{{"_id", oid}}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog by id: %v", err))
	}

	//update the data
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()
	data.AuthorId = blog.GetAuthorId()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse id: %v", err))
	}

	data := &blogItem{}
	filter := bson.D{{"_id", oid}}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog by id: %v", err))
	}
	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.Id.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}
}

func (s *server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := request.GetBlog()
	data := blogItem{
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal errpr: %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Connot convert to oid: %v", err))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

type blogItem struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Title    string             `bson:"content"`
	Content  string             `bson:"title"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errConn := client.Connect(ctx)
	if errConn != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//block until a signal is received
	<-ch
	fmt.Println("Stropping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}
