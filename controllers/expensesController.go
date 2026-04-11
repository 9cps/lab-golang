package controllers

import (
	"net/http"
	"time"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	"github.com/9cps/api-go-gin/services/interfaces"
	"github.com/gin-gonic/gin"
)

type ExpensesController struct {
	expensesServices interfaces.IExpensesServices
}

func NewExpensesController(services interfaces.IExpensesServices) *ExpensesController {
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
//	@Router		/Expenses/CreateExpenses [PUT]
func (c *ExpensesController) CreateExpenses(ctx *gin.Context) {
	req := req_dtos.Expenses{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := c.expensesServices.InsertExpenses(req)
	if result == (models.Expenses{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating expenses"})
		return
	}

	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	})
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
//	@Router		/Expenses/CreateExpensesDetail [PUT]
func (c *ExpensesController) CreateExpensesDetail(ctx *gin.Context) {
	req := req_dtos.ExpensesDetail{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := c.expensesServices.InsertExpensesDetail(req)
	if result == (models.ExpensesDetail{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating expenses detail"})
		return
	}

	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses detail created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	})
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
	listData := c.expensesServices.GetListMoneyCard()
	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData,
	})
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listData := c.expensesServices.GetListMoneyCardDetail(req)
	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData,
	})
}

// UpdateExpensesDetail godoc
//
//	@Summary	Update expenses detail
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		req_dtos.UpdateExpensesDetail	body		req_dtos.UpdateExpensesDetail	true	"UpdateExpensesDetail data"
//
//	@Success	200								{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/UpdateExpensesDetail [PUT]
func (c *ExpensesController) UpdateExpensesDetail(ctx *gin.Context) {
	req := req_dtos.UpdateExpensesDetail{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := c.expensesServices.UpdateExpensesDetail(req)
	if result == (models.ExpensesDetail{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating expenses detail"})
		return
	}

	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses detail updated successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	})
}

// DeleteExpensesDetail godoc
//
//	@Summary	Delete expenses detail
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		req_dtos.DeleteExpensesDetailById	body		req_dtos.DeleteExpensesDetailById	true	"DeleteExpensesDetail data"
//
//	@Success	200									{object}	res_dtos.DefaultResponse
//
//	@Router		/Expenses/DeleteExpensesDetail [DELETE]
func (c *ExpensesController) DeleteExpensesDetail(ctx *gin.Context) {
	req := req_dtos.DeleteExpensesDetailById{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok := c.expensesServices.DeleteExpensesDetail(req); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting expenses detail"})
		return
	}

	ctx.JSON(http.StatusOK, res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Expenses detail deleted successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	})
}
