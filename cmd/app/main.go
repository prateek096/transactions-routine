package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prateek096/transactions-routine/internal/config/db"
	"github.com/prateek096/transactions-routine/internal/handler"
	"github.com/prateek096/transactions-routine/internal/repo"
	"github.com/prateek096/transactions-routine/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading env failed with error : %v", err)
	}

	router := gin.Default()
	dbConn := db.Connect()

	// Initialize Domain
	repo := repo.NewRepo(dbConn)
	accSvc := service.NewAccountService(repo)
	accHandler := handler.NewAccountHandler(accSvc)
	accHandler.RegisterRoutes(router)

	transactionSvc := service.NewTransactionService(repo)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)
	transactionHandler.RegisterRoutes(router)

	// Run router
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
