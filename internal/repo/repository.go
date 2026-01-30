package repo

import (
	"log"

	"github.com/prateek096/transactions-routine/internal/entity"
	"gorm.io/gorm"
)

type Repository interface {
	SaveAccount(account *entity.Account) (*entity.Account, error)
	GetAccountByID(accountID int) (*entity.Account, error)
	SaveTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
}

func NewRepo(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

type RepositoryImpl struct {
	db *gorm.DB
}

func (r *RepositoryImpl) SaveAccount(account *entity.Account) (*entity.Account, error) {
	log.Printf("saving account: %+v", account)
	res := r.db.Create(account)
	if res.Error != nil {
		return nil, res.Error
	}
	return account, nil
}

func (r *RepositoryImpl) GetAccountByID(accountID int) (*entity.Account, error) {
	log.Printf("fetching account with ID: %d", accountID)
	var account entity.Account
	res := r.db.First(&account, accountID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &account, nil
}

func (r *RepositoryImpl) SaveTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	log.Printf("saving transaction: %+v", transaction)
	res := r.db.Create(transaction)
	if res.Error != nil {
		return nil, res.Error
	}
	return transaction, nil
}
