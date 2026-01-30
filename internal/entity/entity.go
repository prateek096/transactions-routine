package entity

import (
	"time"
)

const (
	NormalPurchase      = 1
	InstallmentPurchase = 2
	Withdrawal          = 3
	CreditVoucher       = 4
)

type Account struct {
	AccountId      uint   `json:"account_id" gorm:"primaryKey"`
	DocumentNumber string `json:"document_number" binding:"required" gorm:"uniqueIndex;not null" `
}

type Transaction struct {
	TransactionId   int       `json:"transaction_id" gorm:"primaryKey"`
	AccountId       int       `json:"account_id"  gorm:"uniqueIndex;not null"`
	OperationTypeId int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}
