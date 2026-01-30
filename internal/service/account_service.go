package service

import (
	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/repo"
)

type AccountService interface {
	Create(account *entity.Account) (*entity.Account, error)
	GetById(accountID int) (*entity.Account, error)
}

func NewAccountService(repo repo.Repository) AccountService {
	return &AccountServiceImpl{
		repo: repo,
	}
}

// implementation of AccountService interface
type AccountServiceImpl struct {
	repo repo.Repository
}

func (s *AccountServiceImpl) Create(account *entity.Account) (*entity.Account, error) {
	return s.repo.SaveAccount(account)
}

func (s *AccountServiceImpl) GetById(accountID int) (*entity.Account, error) {
	return s.repo.GetAccountByID(accountID)
}
