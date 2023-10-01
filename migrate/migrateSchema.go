package main

import (
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConncetDatabse()
}

func main() {
	initializers.DB.AutoMigrate(&models.Expenses{}, &models.ExpensesDetail{})
}
