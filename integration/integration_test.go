// server/server_test.go
package server

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/apfelkraepfla/quotes-service/protos/quotespb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testQuoteService struct {
	pb.QuoteServiceServer
}

func (t *testQuoteService) GetQuote(ctx context.Context, req *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	return &pb.QuoteResponse{
		Quote: "Test quote for testing",
	}, nil
}

// startTestServer starts the gRPC server for testing on a random available port.
func startTestServer(t *testing.T) (*grpc.Server, net.Listener) {
	// Get an available port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get an available port: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the QuoteServiceServer
	pb.RegisterQuoteServiceServer(grpcServer, &testQuoteService{})

	// Start the gRPC server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return grpcServer, lis
}

// stopTestServer stops the gRPC server for testing.
func stopTestServer(grpcServer *grpc.Server, lis net.Listener) {
	grpcServer.Stop()
	lis.Close()
}

func TestGetQuote(t *testing.T) {
	// Start the test gRPC server
	grpcServer, lis := startTestServer(t)
	defer stopTestServer(grpcServer, lis)

	// Get the address of the test server
	serverAddr := lis.Addr().String()

	// Create a gRPC client to connect to the test server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial the server: %v", err)
	}
	defer conn.Close()

	// Create a QuoteServiceClient
	client := pb.NewQuoteServiceClient(conn)

	// Perform your gRPC tests using the client...
	req := &pb.QuoteRequest{}
	res, err := client.GetQuote(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to get quote: %v", err)
	}

	expectedQuote := "Test quote for testing"
	if res.Quote != expectedQuote {
		t.Errorf("Expected quote: %s, got: %s", expectedQuote, res.Quote)
	}

	// You can add more test cases as needed...
}
