package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bryanvaz/grpc-gl/protos/go/banking"
	GrpcServer "github.com/bryanvaz/grpc-gl/src/server"
)

var accounts = make(map[string]int32)
var transactions = make(map[string]*banking.Transaction)
var mtx sync.Mutex

var DEBUG = true

func main() {
	s := GrpcServer.NewServer()

	// Start serving incoming connections
	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Print a console message indicating that the server is running
	log.Printf("Server started, listening on port %d", s.Port)
	log.Println("Press Ctrl+C to quit")

	// Block until a signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// Gracefully stop the server
	log.Println("Shutting down server...")
	s.GracefulStop()
}
