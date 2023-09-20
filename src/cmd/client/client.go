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
	last_update := time.Now().UnixMicro()
	if numConns > 1 {
	}

	for i := 0; i < numIter; i++ {
		start := time.Now().UnixMicro()
		_, err := c.Ping(context.Background(), pr)
		if err != nil {
			log.Fatalf("Failed to ping: %v", err)
		}
		end := time.Now().UnixMicro()
		latencies[i] = end - start
		if end-last_update > 500_000 || i == numIter-1 {
			avg_latency := int64(0)
			for li := 0; li <= i; li++ {
				avg_latency += latencies[li]
			}
			avg_latency /= int64(i + 1)
			log.Printf("Iteration %d/%d - Latency: %d μsec (avg)", i+1, numIter, avg_latency)
			last_update = end
		}
	}
	// log.Printf("Response from server: %v", pongResp)

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
