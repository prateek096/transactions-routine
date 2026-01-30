package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/service"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// RegisterRoutes registers the transaction routes
func (h *TransactionHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/transactions")
	{
		group.POST("", h.CreateTransaction)
	}
}

// CreateTransaction handles the creation of a new transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req entity.Transaction

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "desc": err.Error()})
		return
	}

	log.Printf("request to create transaction, AccountId:%v, Amount:%v", req.AccountId, req.Amount)
	transaction, err := h.transactionService.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOperationType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "transaction already exists"})
			return
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "account does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}
