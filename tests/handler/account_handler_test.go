package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/handler"
	"gorm.io/gorm"
)

// Mock service implementing only what's needed
type mockAccountService struct {
	createFn func(*entity.Account) (*entity.Account, error)
	getFn    func(int) (*entity.Account, error)
}

func (m *mockAccountService) Create(a *entity.Account) (*entity.Account, error) { return m.createFn(a) }
func (m *mockAccountService) GetById(id int) (*entity.Account, error)           { return m.getFn(id) }

func setupRouterWithAccountService(mockSvc *mockAccountService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewAccountHandler(mockSvc)
	h.RegisterRoutes(r)
	return r
}

func TestCreateAccount_Success(t *testing.T) {
	mock := &mockAccountService{createFn: func(a *entity.Account) (*entity.Account, error) {
		a.AccountId = 1
		return a, nil
	}}
	r := setupRouterWithAccountService(mock)

	rec := postHelper(r, "/accounts", map[string]string{"document_number": "123"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201; got %d", rec.Code)
	}
	var resp entity.Account
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.AccountId != 1 || resp.DocumentNumber != "123" {
		t.Fatalf("unexpected response body: %+v", resp)
	}
}

func TestCreateAccount_BadRequest(t *testing.T) {
	mock := &mockAccountService{createFn: func(a *entity.Account) (*entity.Account, error) { return nil, nil }}
	r := setupRouterWithAccountService(mock)

	rec := postHelper(r, "/accounts", []byte("not-json"))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400; got %d", rec.Code)
	}
}

func TestCreateAccount_Duplicate(t *testing.T) {
	mock := &mockAccountService{createFn: func(a *entity.Account) (*entity.Account, error) { return nil, gorm.ErrDuplicatedKey }}
	r := setupRouterWithAccountService(mock)

	rec := postHelper(r, "/accounts", map[string]string{"document_number": "123"})

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409; got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["error"] != "account already exists" {
		t.Fatalf("unexpected error message: %v", resp)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	mock := &mockAccountService{getFn: func(int) (*entity.Account, error) { return nil, gorm.ErrRecordNotFound }}
	// Note: handler will translate gorm.ErrRecordNotFound to 404

	r := setupRouterWithAccountService(mock)
	req := httptest.NewRequest(http.MethodGet, "/accounts/999", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404; got %d", rec.Code)
	}
}
