package main

import (
	"context"
	"log"
	"os"

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
	default:
		log.Fatalf("Invalid command provided")
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
