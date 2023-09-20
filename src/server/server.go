package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/bryanvaz/grpc-gl/protos/go/banking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ServerIsRunningError = errors.New("Server is running.")

var accounts = make(map[string]int32)
var transactions = make(map[string]*banking.Transaction)
var mtx sync.Mutex

var DEBUG = true

type Server struct {
	banking.UnimplementedBankingServiceServer
	Port       int
	running    bool
	grpcServer *grpc.Server
}

func NewServer() *Server {
	return &Server{
		Port: 50051,
	}
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()
	banking.RegisterBankingServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	s.grpcServer = grpcServer
	s.running = true
	// Start serving incoming connections
	err = grpcServer.Serve(listener)
	s.running = false
	return err
}

func (s *Server) GracefulStop() {
	if s.running {
		// Gracefully stop the server
		s.grpcServer.GracefulStop()
		s.running = false
	}
}

func (s *Server) Ping(ctx context.Context, req *banking.PingRequest) (*banking.PingResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	return &banking.PingResponse{Message: "Pong"}, nil
}

func (s *Server) MakeTransaction(ctx context.Context, req *banking.TransactionRequest) (*banking.TransactionResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	fromBalance, ok1 := accounts[req.FromAccountId]
	toBalance, ok2 := accounts[req.ToAccountId]

	if !ok1 || !ok2 {
		return &banking.TransactionResponse{Success: false, Message: "Account not found"}, nil
	}

	// if fromBalance < req.Amount {
	// 	return &banking.TransactionResponse{Success: false, Message: "Insufficient balance"}, nil
	// }

	transactionID := fmt.Sprintf("%d", time.Now().UnixNano())

	accounts[req.FromAccountId] -= req.Amount
	accounts[req.ToAccountId] += req.Amount

	transaction := &banking.Transaction{
		TransactionId: transactionID,
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	transactions[transactionID] = transaction

	if DEBUG {
		log.Printf(
			"MakeTransaction: Old balances: From: %d (%s), To: %d (%s)\n",
			fromBalance, req.FromAccountId, toBalance, req.ToAccountId,
		)
		log.Printf(
			"MakeTransaction: ID: %s, From: %s, To: %s, Amount: %d\n",
			transactionID, req.FromAccountId, req.ToAccountId, req.Amount,
		)
		log.Printf(
			"MakeTransaction: New balances: From: %d (%s), To: %d (%s)\n",
			accounts[req.FromAccountId], req.FromAccountId, accounts[req.ToAccountId], req.ToAccountId,
		)
	}

	return &banking.TransactionResponse{TransactionId: transactionID, Success: true, Message: "Transaction Successful"}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *banking.BalanceRequest) (*banking.BalanceResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	balance, ok := accounts[req.AccountId]
	if !ok {
		return nil, fmt.Errorf("Account not found")
	}

	if DEBUG {
		log.Println("GetBalance: ID:", req.AccountId, "Balance:", balance)
	}

	return &banking.BalanceResponse{Balance: balance}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *banking.AccountRequest) (*banking.AccountResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	accountID := fmt.Sprintf("%d", time.Now().UnixMilli())
	accounts[accountID] = req.InitialBalance

	if DEBUG {
		log.Println("GetBalance: ID:", accountID, "Balance:", req.InitialBalance)
	}

	return &banking.AccountResponse{AccountId: accountID}, nil
}

func (s *Server) GetTransactionDetails(ctx context.Context, req *banking.TransactionDetailsRequest) (*banking.TransactionDetailsResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	transaction, ok := transactions[req.TransactionId]
	if !ok {
		return nil, fmt.Errorf("Transaction not found")
	}

	if DEBUG {
		log.Println("GetTransactionDetails: ID:", req.TransactionId, "Transaction:", transaction)
	}

	return &banking.TransactionDetailsResponse{Transaction: transaction}, nil
}

func (s *Server) ListAccount(ctx context.Context, req *banking.ListAccountRequest) (*banking.ListAccountResponse, error) {
	mtx.Lock()
	defer mtx.Unlock()

	var accountList []*banking.Account

	for id, balance := range accounts {
		accountList = append(accountList, &banking.Account{Id: id, Balance: balance})
	}

	if DEBUG {
		log.Println("ListAccount: Accounts:", accountList)
	}

	return &banking.ListAccountResponse{Accounts: accountList}, nil
}
