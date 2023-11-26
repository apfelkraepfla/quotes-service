package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/apfelkraepfla/quotes_service/server"
	"github.com/joho/godotenv"
)

var port = flag.Int("port", getEnvAsInt("MYAPP_PORT", 0), "The server port")

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
	flag.Parse()

	// Use the specified port or the one from the environment variable or a default value
	config := server.ServerConfig{
		Port: *port,
	}

	server.StartServer(config)
}
