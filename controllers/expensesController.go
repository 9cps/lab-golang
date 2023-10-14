package controllers

import (
	"time"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/helper"
	"github.com/9cps/api-go-gin/models"
	services "github.com/9cps/api-go-gin/services"
	"github.com/gin-gonic/gin"
)

type ExpensesController struct {
	expensesServices services.ExpensesServices
}

func NewExpensesController(services services.ExpensesServices) *ExpensesController {
	return &ExpensesController{
		expensesServices: services,
	}
}

// CreateExpenses godoc
//
//	@Summary	Create expenses
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		req_dtos.Expenses body req_dtos.Expenses true "Expense data"
//
//	@Success	200	{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/create [POST]
func (c *ExpensesController) CreateExpenses(ctx *gin.Context) {

	req := req_dtos.Expenses{}
	err := ctx.ShouldBindJSON(&req)
	helper.ErrorPanic(err)

	result := c.expensesServices.InsertExpenses(req)
	if result == (models.Expenses{}) {
		ctx.JSON(400, gin.H{
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

	ctx.JSON(200, response)
}

// func CreateExpensesDetail(c *gin.Context) {
// 	result := services.InsertExpensesDetail(c)
// 	if result == (models.ExpensesDetail{}) {
// 		c.JSON(400, gin.H{
// 			"error": "Error creating expenses detail",
// 		})
// 		return
// 	}
// 	response := res_dtos.DefaultResponse{
// 		Status:  string(res_dtos.Success),
// 		Message: "Expenses detail created successfully",
// 		Date:    time.Now().Format("02/01/2006 15:04:05"),
// 		Data:    result,
// 	}

// 	c.JSON(200, response)
// }

// func GetListMoneyCard(c *gin.Context) {
// 	// Find All MoneyCard
// 	listData := services.GetListMoneyCard(c)
// 	response := res_dtos.DefaultResponse{
// 		Status: string(res_dtos.Success),
// 		Date:   time.Now().Format("02/01/2006 15:04:05"),
// 		Data:   listData, // Data retrieved from the SQL query
// 	}
// 	// Return data
// 	c.JSON(200, response)
// }

// func GetListMoneyCardDetail(c *gin.Context) {
// 	// Find All MoneyCard
// 	listData := services.GetListMoneyCardDetail(c)
// 	response := res_dtos.DefaultResponse{
// 		Status: string(res_dtos.Success),
// 		Date:   time.Now().Format("02/01/2006 15:04:05"),
// 		Data:   listData, // Data retrieved from the SQL query
// 	}
// 	// Return data
// 	c.JSON(200, response)
// }
