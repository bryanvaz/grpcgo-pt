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

// func main_old() {
// 	port := 50051
// 	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	banking.RegisterBankingServiceServer(grpcServer, &server{})

// 	reflection.Register(grpcServer)

// 	// Start serving incoming connections
// 	go func() {
// 		if err := grpcServer.Serve(listener); err != nil {
// 			log.Fatalf("Failed to serve: %v", err)
// 		}
// 	}()

// 	// Print a console message indicating that the server is running
// 	log.Printf("Server started, listening on port %d", port)
// 	log.Println("Press Ctrl+C to quit")

// 	// Block until a signal is received
// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
// 	<-signalChan

// 	// Gracefully stop the server
// 	log.Println("Shutting down server...")
// 	grpcServer.GracefulStop()

// if err := grpcServer.Serve(listener); err != nil {
// 	log.Fatalf("Failed to serve: %v", err)
// }
// }

// type server struct {
// 	banking.UnimplementedBankingServiceServer
// }

// func (s *server) MakeTransaction(ctx context.Context, req *banking.TransactionRequest) (*banking.TransactionResponse, error) {
// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	fromBalance, ok1 := accounts[req.FromAccountId]
// 	toBalance, ok2 := accounts[req.ToAccountId]

// 	if !ok1 || !ok2 {
// 		return &banking.TransactionResponse{Success: false, Message: "Account not found"}, nil
// 	}

// 	// if fromBalance < req.Amount {
// 	// 	return &banking.TransactionResponse{Success: false, Message: "Insufficient balance"}, nil
// 	// }

// 	transactionID := fmt.Sprintf("%d", time.Now().UnixNano())

// 	accounts[req.FromAccountId] -= req.Amount
// 	accounts[req.ToAccountId] += req.Amount

// 	transaction := &banking.Transaction{
// 		TransactionId: transactionID,
// 		FromAccountId: req.FromAccountId,
// 		ToAccountId:   req.ToAccountId,
// 		Amount:        req.Amount,
// 	}

// 	transactions[transactionID] = transaction

// 	if DEBUG {
// 		log.Printf(
// 			"MakeTransaction: Old balances: From: %d (%s), To: %d (%s)\n",
// 			fromBalance, req.FromAccountId, toBalance, req.ToAccountId,
// 		)
// 		log.Printf(
// 			"MakeTransaction: ID: %s, From: %s, To: %s, Amount: %d\n",
// 			transactionID, req.FromAccountId, req.ToAccountId, req.Amount,
// 		)
// 		log.Printf(
// 			"MakeTransaction: New balances: From: %d (%s), To: %d (%s)\n",
// 			accounts[req.FromAccountId], req.FromAccountId, accounts[req.ToAccountId], req.ToAccountId,
// 		)
// 	}

// 	return &banking.TransactionResponse{TransactionId: transactionID, Success: true, Message: "Transaction Successful"}, nil
// }

// func (s *server) GetBalance(ctx context.Context, req *banking.BalanceRequest) (*banking.BalanceResponse, error) {
// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	balance, ok := accounts[req.AccountId]
// 	if !ok {
// 		return nil, fmt.Errorf("Account not found")
// 	}

// 	if DEBUG {
// 		log.Println("GetBalance: ID:", req.AccountId, "Balance:", balance)
// 	}

// 	return &banking.BalanceResponse{Balance: balance}, nil
// }

// func (s *server) CreateAccount(ctx context.Context, req *banking.AccountRequest) (*banking.AccountResponse, error) {
// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	accountID := fmt.Sprintf("%d", time.Now().UnixMilli())
// 	accounts[accountID] = req.InitialBalance

// 	if DEBUG {
// 		log.Println("GetBalance: ID:", accountID, "Balance:", req.InitialBalance)
// 	}

// 	return &banking.AccountResponse{AccountId: accountID}, nil
// }

// func (s *server) GetTransactionDetails(ctx context.Context, req *banking.TransactionDetailsRequest) (*banking.TransactionDetailsResponse, error) {
// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	transaction, ok := transactions[req.TransactionId]
// 	if !ok {
// 		return nil, fmt.Errorf("Transaction not found")
// 	}

// 	if DEBUG {
// 		log.Println("GetTransactionDetails: ID:", req.TransactionId, "Transaction:", transaction)
// 	}

// 	return &banking.TransactionDetailsResponse{Transaction: transaction}, nil
// }

// func (s *server) ListAccount(ctx context.Context, req *banking.ListAccountRequest) (*banking.ListAccountResponse, error) {
// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	var accountList []*banking.Account

// 	for id, balance := range accounts {
// 		accountList = append(accountList, &banking.Account{Id: id, Balance: balance})
// 	}

// 	if DEBUG {
// 		log.Println("ListAccount: Accounts:", accountList)
// 	}

// 	return &banking.ListAccountResponse{Accounts: accountList}, nil
// }
