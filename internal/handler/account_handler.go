package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prateek096/transactions-routine/internal/entity"
	"github.com/prateek096/transactions-routine/internal/service"
	"gorm.io/gorm"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// RegisterRoutes registers the account routes
func (h *AccountHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/accounts")
	{
		group.POST("", h.CreateAccount)
		group.GET("/:accountId", h.GetAccountByID)
	}
}

// CreateAccount handles the creation of a new account
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req entity.Account

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "desc": err.Error()})
		return
	}

	log.Printf("request to create account, Id:%v", req)
	account, err := h.accountService.Create(&req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "account already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccountByID handles fetching an account by its ID
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	accountId, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid account id",
			"desc":  "id must be a number",
		})
		return
	}
	log.Printf("request to get account, Id:%d", accountId)
	account, err := h.accountService.GetById(accountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, account)
}
