package controllers

import (
	"time"

	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	services "github.com/9cps/api-go-gin/services"
	"github.com/gin-gonic/gin"
)

func CreateExpenses(c *gin.Context) {
	result := services.InsertExpenses(c)
	if result == (models.Expenses{}) {
		c.JSON(400, gin.H{
			"error": "Error creating expenses",
		})
		return
	}
	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	}

	c.JSON(200, response)
}

func CreateExpensesDetail(c *gin.Context) {
	result := services.InsertExpensesDetail(c)
	if result == (models.ExpensesDetail{}) {
		c.JSON(400, gin.H{
			"error": "Error creating expenses detail",
		})
		return
	}
	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses detail created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	}

	c.JSON(200, response)
}

func GetListMoneyCard(c *gin.Context) {
	// Find All MoneyCard
	listData := services.GetListMoneyCard(c)
	response := res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData, // Data retrieved from the SQL query
	}
	// Return data
	c.JSON(200, response)
}

// func DeleteFriend(c *gin.Context) {
// 	// Find friend by keyword
// 	friendData := services.DeleteFriend(c)
// 	response := res_dtos.DefaultResponse{
// 		Status: string(res_dtos.Success),
// 		Date:   time.Now().Format("02/01/2006 15:04:05"),
// 		Data:   friendData, // Data retrieved from the SQL query
// 	}

// 	// Return data
// 	c.JSON(200, response)
// }
