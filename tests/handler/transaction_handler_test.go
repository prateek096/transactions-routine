package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/handler"
	"github.com/prateek096/transactions-routine/internal/service"
	"gorm.io/gorm"
)

// Mock transaction service implementing only Create
type mockTransactionService struct {
	createFn func(*entity.Transaction) (*entity.Transaction, error)
}

func (m *mockTransactionService) Create(t *entity.Transaction) (*entity.Transaction, error) {
	return m.createFn(t)
}

func setupRouterWithTransactionService(mockSvc *mockTransactionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewTransactionHandler(mockSvc)
	h.RegisterRoutes(r)
	return r
}

// helper to POST JSON and return recorder
func postHelper(r *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestCreateTransaction_Success(t *testing.T) {
	mock := &mockTransactionService{createFn: func(tx *entity.Transaction) (*entity.Transaction, error) {
		tx.TransactionId = 10
		return tx, nil
	}}
	r := setupRouterWithTransactionService(mock)

	rec := postHelper(r, "/transactions", map[string]interface{}{"account_id": 1, "operation_type_id": 1, "amount": 100.0})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201; got %d", rec.Code)
	}
	var resp entity.Transaction
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.TransactionId != 10 {
		t.Fatalf("unexpected transaction id in response: %v", resp.TransactionId)
	}
}

func TestCreateTransaction_BadRequest(t *testing.T) {
	mock := &mockTransactionService{createFn: func(tx *entity.Transaction) (*entity.Transaction, error) { return nil, nil }}
	r := setupRouterWithTransactionService(mock)

	rec := postHelper(r, "/transactions", []byte("not-json"))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400; got %d", rec.Code)
	}
}

func TestCreateTransaction_InvalidOperation(t *testing.T) {
	mock := &mockTransactionService{createFn: func(tx *entity.Transaction) (*entity.Transaction, error) {
		return nil, service.ErrInvalidOperationType
	}}
	r := setupRouterWithTransactionService(mock)
	rec := postHelper(r, "/transactions", map[string]interface{}{"account_id": 1, "operation_type_id": 99, "amount": 10.0})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400; got %d", rec.Code)
	}
}

func TestCreateTransaction_Duplicate(t *testing.T) {
	mock := &mockTransactionService{createFn: func(tx *entity.Transaction) (*entity.Transaction, error) {
		return nil, gorm.ErrDuplicatedKey
	}}
	r := setupRouterWithTransactionService(mock)

	rec := postHelper(r, "/transactions", map[string]interface{}{"account_id": 1, "operation_type_id": 1, "amount": 10.0})

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409; got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["error"] != "transaction already exists" {
		t.Fatalf("unexpected error message: %v", resp)
	}
}
