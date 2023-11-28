package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/apfelkraepfla/quotes-service/protos/quotespb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type quoteServer struct {
	pb.QuoteServiceServer
}

type ServerConfig struct {
	Port int
}

func (s *quoteServer) GetQuote(ctx context.Context, req *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	// Return a hardcoded philosophical quote
	return &pb.QuoteResponse{
		Quote: "The only true wisdom is in knowing you know nothing. - Socrates",
	}, nil
}

func StartServer(config ServerConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("Failed to listen at %v: %v", config.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQuoteServiceServer(grpcServer, &quoteServer{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %d", config.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
