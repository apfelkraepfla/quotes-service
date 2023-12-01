package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/apfelkraepfla/quotes_service/server"
	"github.com/joho/godotenv"
)

var rpcPort = flag.Int("rpc_port", getEnvAsInt("RPC_PORT", 0), "The rpc server port")
var port = flag.Int("port", getEnvAsInt("PORT", 0), "The REST API server port")
var appEnv = os.Getenv("APP_ENV")

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		// Load environment variables from .env file
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file")
		}
		// Try to get the environment variable again
		value, _ = os.LookupEnv(key)
	}

	// Parse the value to an integer or use the default
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if appEnv == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	logger.Info("First slog logging message")

	flag.Parse()

	// Use the specified port or the one from the environment variable or a default value
	config := server.ServerConfig{
		RpcPort: *rpcPort,
		Port:    *port,
	}

	go server.StartServer(config, *logger)
	go server.StartRPCGatewayServer(config, *logger)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	logger.Info("Termination signal received. Exiting...")
}
