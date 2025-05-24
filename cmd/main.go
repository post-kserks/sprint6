package main

import (
	"log"
	"os"

	"morse-converter/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	srv := server.New(logger)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
