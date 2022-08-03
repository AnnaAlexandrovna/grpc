package main

import (
	"books/bookshop/pb"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var SampleBooks = []*pb.Book{
	{
		Id:        "1",
		Title:     "The Hitchhiker's Guide to the Galaxy",
		Author:    "Douglas Adams",
		PageCount: 42,
	},
	{
		Id:        "2",
		Title:     "The Lord of the Rings",
		Author:    "J.R.R. Tolkien",
		PageCount: 1234,
	},
}

type server struct {
	pb.UnimplementedInventoryServer
}

func (s *server) GetBookList(ctx context.Context, in *pb.GetBookListRequest) (*pb.GetBookListResponse, error) {
	return &pb.GetBookListResponse{
		Books: getSampleBooks(),
	}, nil
}

func (s *server) GetBookById(ctx context.Context, in *pb.GetBookByIdRequest) (*pb.Book, error) {
	book := getBookById(in.Id)
	if book != nil {
		return getBookById(in.Id), nil
	}
	err := status.Error(codes.NotFound, "id was not found")
	return nil, err
}
func getSampleBooks() []*pb.Book {
	return SampleBooks
}
func getBookById(id string) *pb.Book {
	for _, el := range SampleBooks {
		if el.Id == id {
			return el
		}
	}
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterInventoryServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
