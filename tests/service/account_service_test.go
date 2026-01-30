package service_test

import (
	"reflect"
	"testing"

	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/service"
)

type mockRepoAcc struct {
	savedAccount *entity.Account
	getAccount   *entity.Account
	getErr       error
	saveErr      error
}

func (m *mockRepoAcc) SaveAccount(account *entity.Account) (*entity.Account, error) {
	m.savedAccount = account
	return account, m.saveErr
}

func (m *mockRepoAcc) GetAccountByID(accountID int) (*entity.Account, error) {
	return m.getAccount, m.getErr
}

func (m *mockRepoAcc) SaveTransaction(tx *entity.Transaction) (*entity.Transaction, error) {
	return nil, nil
}

func TestAccountService_Create(t *testing.T) {
	mock := &mockRepoAcc{}
	svc := service.NewAccountService(mock)

	acc := &entity.Account{DocumentNumber: "12345678901"}
	got, err := svc.Create(acc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, acc) {
		t.Fatalf("expected saved account to be returned; got %v, want %v", got, acc)
	}
	if mock.savedAccount == nil || mock.savedAccount.DocumentNumber != acc.DocumentNumber {
		t.Fatal("repo.SaveAccount was not called")
	}
}

func TestAccountService_GetById(t *testing.T) {
	expected := &entity.Account{AccountId: 10, DocumentNumber: "ABC"}
	mock := &mockRepoAcc{getAccount: expected}
	svc := service.NewAccountService(mock)

	got, err := svc.GetById(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected account expected; got %v, want %v", got, expected)
	}
}
