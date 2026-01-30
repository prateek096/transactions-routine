package service

import (
	"errors"
	"time"

	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/repo"
)

type TransactionService interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
}

func NewTransactionService(repo repo.Repository) TransactionService {
	return &TransactionServiceImpl{repo: repo}
}

// implementation of TransactionService interface
type TransactionServiceImpl struct {
	repo repo.Repository
}

func (s *TransactionServiceImpl) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	if err := s.validateTransaction(transaction); err != nil {
		return nil, err
	}
	s.processTransaction(transaction)
	transaction.EventDate = time.Now()
	return s.repo.SaveTransaction(transaction)
}

func (s *TransactionServiceImpl) processTransaction(transaction *entity.Transaction) {
	// For purchase operations, amount should be negative
	if transaction.OperationTypeId == entity.NormalPurchase || transaction.OperationTypeId == entity.InstallmentPurchase || transaction.OperationTypeId == entity.Withdrawal {
		if transaction.Amount > 0 {
			transaction.Amount = -transaction.Amount
		}
	}
}

var ErrInvalidOperationType = errors.New("invalid operation type")

func (s *TransactionServiceImpl) validateTransaction(transaction *entity.Transaction) error {
	if transaction.OperationTypeId < 1 || transaction.OperationTypeId > 4 {
		return ErrInvalidOperationType
	}
	return nil
}
