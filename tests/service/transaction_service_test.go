package service_test

import (
	"errors"
	"testing"

	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/service"
)

type mockRepoTx struct {
	savedTx *entity.Transaction
	saveErr error
}

func (m *mockRepoTx) SaveAccount(account *entity.Account) (*entity.Account, error) { return nil, nil }
func (m *mockRepoTx) GetAccountByID(accountID int) (*entity.Account, error)        { return nil, nil }
func (m *mockRepoTx) SaveTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	if m.saveErr != nil {
		return nil, m.saveErr
	}
	m.savedTx = transaction
	transaction.TransactionId = 123
	return transaction, nil
}

func TestTransactionService_PurchaseAmountNegative(t *testing.T) {
	mock := &mockRepoTx{}
	svc := service.NewTransactionService(mock)

	tx := &entity.Transaction{AccountId: 1, OperationTypeId: entity.NormalPurchase, Amount: 100.0}
	got, err := svc.Create(tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != -100.0 {
		t.Fatalf("expected amount to be negative for purchase; got %v", got.Amount)
	}
	if mock.savedTx != got {
		t.Fatalf("expected repo.SaveTransaction to be called")
	}
}

func TestTransactionService_CreditVoucherPositive(t *testing.T) {
	mock := &mockRepoTx{}
	svc := service.NewTransactionService(mock)

	tx := &entity.Transaction{AccountId: 1, OperationTypeId: entity.CreditVoucher, Amount: 150.0}
	got, err := svc.Create(tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Amount != 150.0 {
		t.Fatalf("expected amount to keep sign for credit voucher; got %v", got.Amount)
	}
}

func TestTransactionService_InvalidOperationType(t *testing.T) {
	mock := &mockRepoTx{}
	svc := service.NewTransactionService(mock)

	tx := &entity.Transaction{AccountId: 1, OperationTypeId: 99, Amount: 10.0}
	_, err := svc.Create(tx)
	if err == nil {
		t.Fatalf("expected error for invalid operation type; got nil")
	}
	if !errors.Is(err, service.ErrInvalidOperationType) {
		t.Fatalf("expected ErrInvalidOperationType; got %v", err)
	}
}
