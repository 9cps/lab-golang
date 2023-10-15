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
//	@Param		req_dtos.Expenses	body		req_dtos.Expenses	true	"Expense data"
//
//	@Success	200					{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/CreateExpenses [POST]
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

// CreateExpensesDetail godoc
//
//	@Summary	Create expenses
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		req_dtos.ExpensesDetail	body		req_dtos.ExpensesDetail	true	"ExpensesDetail data"
//
//	@Success	200						{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/CreateExpensesDetail [POST]
func (c *ExpensesController) CreateExpensesDetail(ctx *gin.Context) {
	req := req_dtos.ExpensesDetail{}
	err := ctx.ShouldBindJSON(&req)
	helper.ErrorPanic(err)

	result := c.expensesServices.InsertExpensesDetail(req)
	if result == (models.ExpensesDetail{}) {
		ctx.JSON(400, gin.H{
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

	ctx.JSON(200, response)
}

// GetListMoneyCard godoc
//
//	@Summary	Get list money item
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//
//	@Success	200	{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/GetListMoneyCard [GET]
func (c *ExpensesController) GetListMoneyCard(ctx *gin.Context) {

	// Find All MoneyCard
	listData := c.expensesServices.GetListMoneyCard()
	response := res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData, // Data retrieved from the SQL query
	}
	// Return data
	ctx.JSON(200, response)
}

// GetListMoneyCardDetail godoc
//
//	@Summary	Get list money detail item
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		req_dtos.GetExpensesDetailById	body		req_dtos.GetExpensesDetailById	true	"GetListMoneyCardDetail data"
//
//	@Success	200								{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/GetListMoneyCardDetail [POST]
func (c *ExpensesController) GetListMoneyCardDetail(ctx *gin.Context) {
	req := req_dtos.GetExpensesDetailById{}
	err := ctx.ShouldBindJSON(&req)
	helper.ErrorPanic(err)

	// Find All MoneyCard
	listData := c.expensesServices.GetListMoneyCardDetail(req)
	response := res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData, // Data retrieved from the SQL query
	}
	// Return data
	ctx.JSON(200, response)
}
