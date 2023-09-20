package server

import (
	"context"
	"testing"

	"github.com/bryanvaz/grpc-gl/protos/go/banking"
	"github.com/stretchr/testify/assert"
)

func getNewTestServer() *Server {
	s := NewServer()
	s.TestMode(true)
	return s
}

func TestServer_Ping(t *testing.T) {
	s := getNewTestServer()
	req := &banking.PingRequest{}
	expected := &banking.PingResponse{Message: "Pong"}

	res, err := s.Ping(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestServer_MakeTransaction(t *testing.T) {
	s := getNewTestServer()
	ca1, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})
	ca2, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})

	req := &banking.TransactionRequest{
		FromAccountId: ca1.AccountId,
		ToAccountId:   ca2.AccountId,
		Amount:        50,
	}
	expected := &banking.TransactionResponse{
		TransactionId: ".*",
		Success:       true,
		Message:       "Transaction Successful",
	}

	res, err := s.MakeTransaction(context.Background(), req)

	assert.NoError(t, err)
	assert.Regexp(t, expected.TransactionId, res.TransactionId)
	assert.Equal(t, expected.Success, res.Success)
	assert.Equal(t, expected.Message, res.Message)
}

func TestServer_GetBalance(t *testing.T) {
	s := getNewTestServer()
	ca1, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})

	req := &banking.BalanceRequest{AccountId: ca1.AccountId}
	expected := &banking.BalanceResponse{Balance: 100}

	res, err := s.GetBalance(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestServer_CreateAccount(t *testing.T) {
	s := getNewTestServer()
	req := &banking.AccountRequest{InitialBalance: 100}
	expected := &banking.AccountResponse{AccountId: ".*"}

	res, err := s.CreateAccount(context.Background(), req)

	assert.NoError(t, err)
	assert.Regexp(t, expected.AccountId, res.AccountId)
}

func TestServer_GetTransactionDetails(t *testing.T) {
	s := getNewTestServer()
	ca1, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})
	ca2, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})
	tx1, _ := s.MakeTransaction(context.Background(), &banking.TransactionRequest{
		FromAccountId: ca1.AccountId,
		ToAccountId:   ca2.AccountId,
		Amount:        50,
	})

	req := &banking.TransactionDetailsRequest{TransactionId: tx1.TransactionId}
	expected := &banking.TransactionDetailsResponse{
		Transaction: &banking.Transaction{
			TransactionId: tx1.TransactionId,
			FromAccountId: ca1.AccountId,
			ToAccountId:   ca2.AccountId,
			Amount:        50,
		},
	}

	res, err := s.GetTransactionDetails(context.Background(), req)

	assert.NoError(t, err)
	assert.Regexp(t, expected.Transaction.TransactionId, res.Transaction.TransactionId)
	assert.Equal(t, expected.Transaction.FromAccountId, res.Transaction.FromAccountId)
	assert.Equal(t, expected.Transaction.ToAccountId, res.Transaction.ToAccountId)
	assert.Equal(t, expected.Transaction.Amount, res.Transaction.Amount)
}

func TestServer_ListAccount(t *testing.T) {
	s := getNewTestServer()
	ca1, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})
	ca2, _ := s.CreateAccount(context.Background(), &banking.AccountRequest{InitialBalance: 100})

	req := &banking.ListAccountRequest{}
	expected := &banking.ListAccountResponse{
		Accounts: []*banking.Account{
			{Id: ca1.AccountId, Balance: 100},
			{Id: ca2.AccountId, Balance: 100},
		},
	}

	res, err := s.ListAccount(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, len(expected.Accounts), len(res.Accounts))
	assert.Equal(t, expected, res)
}
