// server/server_test.go
package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"

	pb "github.com/siarener/quotes-service/protos/quotespb"
	"github.com/siarener/quotes_service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestQuoteService_GetQuote(t *testing.T) {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	RpcPort := 3032
	config := server.ServerConfig{RpcPort: RpcPort, Port: 3033}

	// Start the test gRPC server
	go server.StartServer(config, *logger)

	stopChan := make(chan os.Signal, 1)

	// Create a gRPC client to connect to the test server
	conn, err := grpc.Dial(fmt.Sprintf(":%d", RpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial the server: %v", err)
	}
	defer conn.Close()

	// Create a QuoteServiceClient
	client := pb.NewQuoteServiceClient(conn)

	// Perform your gRPC tests using the client...
	getReq := &pb.QuoteRequest{}
	getRes, err := client.GetQuote(context.Background(), getReq)
	if err != nil {
		t.Fatalf("Failed to get quote: %v", err)
	}

	expectedQuote := "Test quote for testing"
	if getRes.Quote != expectedQuote {
		t.Errorf("Expected quote: %s, got: %s", expectedQuote, getRes.Quote)
	}
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
}

func TestQuoteService_StoreQuote(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	RpcPort := 3032
	config := server.ServerConfig{RpcPort: RpcPort, Port: 3033}

	// Start the test gRPC server
	go server.StartServer(config, *logger)

	stopChan := make(chan os.Signal, 1)

	// Create a gRPC client to connect to the test server
	conn, err := grpc.Dial(fmt.Sprintf(":%d", RpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial the server: %v", err)
	}
	defer conn.Close()

	// Create a QuoteServiceClient
	client := pb.NewQuoteServiceClient(conn)

	// Perform your gRPC tests using the client...
	storeReq := &pb.StoreQuoteRequest{Quote: "New test quote"}
	storeRes, err := client.StoreQuote(context.Background(), storeReq)
	log.Printf("Here:")
	if err != nil {
		t.Fatalf("Failed to store quote: %v", err)
	}
	log.Printf("Here: %v", storeRes)
	//expectedRes := &emptypb.Empty{}
	//if storeRes != expectedRes {
	//	t.Errorf("Expected response: %s, got: %s", expectedRes, storeRes)
	//}
	// You can add more test cases as needed...
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
}
