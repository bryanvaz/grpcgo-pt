package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/bryanvaz/grpc-gl/protos/go/banking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalf("No command provided")
	}
	action := args[0]

	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client for the banking service
	client := pb.NewBankingServiceClient(conn)

	switch action {
	case "create":
		createAccount(client)
	case "balance":
		getBalance(client)
	case "plow":
		plow(client)
	default:
		log.Fatalf("Invalid command provided")
	}
}

func plow(c pb.BankingServiceClient) {
	numIter := 1
	numConns := 1

	if len(os.Args) > 2 {
		for pos, value := range os.Args {
			if value == "-n" && pos < len(os.Args)-1 {
				numIter, _ = strconv.Atoi(os.Args[pos+1])
			}
		}
		for pos, value := range os.Args {
			if value == "-c" && pos < len(os.Args)-1 {
				numConns, _ = strconv.Atoi(os.Args[pos+1])
			}
		}
	}

	pr := &pb.PingRequest{
		Message: "ping",
	}
	latencies := make([]int64, numIter)

	jobs := make(chan int, numIter)
	results := make(chan int64, numIter)

	// Start the workers
	for i := 0; i < numConns; i++ {
		go func() {
			for j := range jobs {
				_ = j
				start := time.Now().UnixMicro()
				_, err := c.Ping(context.Background(), pr)
				if err != nil {
					log.Fatalf("Failed to ping: %v", err)
				}
				end := time.Now().UnixMicro()
				results <- end - start
			}
		}()
	}

	// Send the jobs
	for i := 0; i < numIter; i++ {
		jobs <- i
	}
	close(jobs)

	// Collect the results
	start_time := time.Now().UnixMicro()
	last_update := time.Now().UnixMicro()
	for i := 0; i < numIter; i++ {
		latencies[i] = <-results
		if last_update < time.Now().UnixMicro()-500_000 || i == numIter-1 {
			avg_latency := int64(0)
			for li := 0; li <= i; li++ {
				avg_latency += latencies[li]
			}
			avg_latency /= int64(i + 1)
			last_update = time.Now().UnixMicro()
			elapsed := last_update - start_time
			pct_complete := float64(i+1) / float64(numIter) * 100.0
			req_per_sec := float64(i+1) / float64(elapsed) * 1_000_000.0
			log.Printf(
				"Iteration %d/%d (%.2f%%)- %.1f sec elapsed - Latency: %d μsec (avg) - RPS: %.1f",
				i+1, numIter, pct_complete, float64(elapsed)/1_000_000.0, avg_latency, req_per_sec,
			)
		}
	}
}

func createAccount(c pb.BankingServiceClient) {
	// Read CLI arguments into a string array
	// args := os.Args[1:]
	// fmt.Println(args)
	// Call the CreateAccount method
	ar := &pb.AccountRequest{
		InitialBalance: 50000.0,
	}
	createAccountResponse, err := c.CreateAccount(context.Background(), ar)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}
	log.Printf("Response from server: %v", createAccountResponse)
}

func getBalance(c pb.BankingServiceClient) {
	// Read CLI arguments into a string array
	args := os.Args[2:]
	if len(args) == 0 {
		log.Fatalf("No account ID provided")
	}
	accountID := args[0]
	getBalanceRequest := &pb.BalanceRequest{
		AccountId: accountID,
	}
	getBalanceResponse, err := c.GetBalance(context.Background(), getBalanceRequest)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	log.Printf("Balance for account %s: %v", accountID, getBalanceResponse.Balance)
}
