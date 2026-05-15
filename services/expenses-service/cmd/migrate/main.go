package main

import (
	"log"

	"github.com/9cps/api-go-gin/services/expenses-service/internal/config"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&model.Expenses{}, &model.ExpensesDetail{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	log.Println("migration completed successfully")
}
