package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/siarener/quotes-service/protos/quotespb"
	pb "github.com/siarener/quotes-service/protos/quotespb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type quoteServer struct {
	pb.QuoteServiceServer
}

type ServerConfig struct {
	RpcPort int
	Port    int
}

func (s *quoteServer) GetQuote(ctx context.Context, req *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	// Return a hardcoded philosophical quote
	return &pb.QuoteResponse{
		Quote: "The only true wisdom is in knowing you know nothing. - Socrates",
	}, nil
}

func (s *quoteServer) StoreQuote(ctx context.Context, req *pb.StoreQuoteRequest) (*pb.Empty, error) {
	// Store a quote
	log.Printf("Storing the quote: %v", req.Quote)
	return &pb.Empty{}, nil
}

func StartServer(config ServerConfig, logger slog.Logger) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.RpcPort))
	if err != nil {
		logger.Error("Failed to listen to port", slog.Int("port", config.RpcPort), err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQuoteServiceServer(grpcServer, &quoteServer{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	logger.Info("gRPC server is running", slog.Int("port", config.RpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve", err)
	}
}

func StartRPCGatewayServer(config ServerConfig, logger slog.Logger) {
	gwmux := runtime.NewServeMux()
	err := quotespb.RegisterQuoteServiceHandlerFromEndpoint(
		context.Background(),
		gwmux, ":"+fmt.Sprintf(":%d", config.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	if err != nil {
		log.Fatal(err)
	}
	gwServer := &http.Server{
		Addr:    ":" + fmt.Sprintf(":%d", config.Port),
		Handler: gwmux,
	}

	logger.Info("Serving gRPC-Gateway on: ", slog.Int("port", config.Port))
	log.Fatalln(gwServer.ListenAndServe())
}
