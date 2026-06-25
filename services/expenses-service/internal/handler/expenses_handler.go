package handler

import (
	"net/http"
	"time"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type ExpensesHandler interface {
	CreateExpenses(ctx *gin.Context)
	GetListMoneyCard(ctx *gin.Context)
	GetListMoneyCardDetail(ctx *gin.Context)
	UpsertExpensesDetail(ctx *gin.Context)
	DeleteExpensesDetail(ctx *gin.Context)
}

type expensesHandler struct {
	expensesService interfaces.ExpensesService
}

func NewExpensesHandler(svc interfaces.ExpensesService) ExpensesHandler {
	return &expensesHandler{expensesService: svc}
}

// CreateExpenses godoc
//
//	@Summary	Create a monthly expense card
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.Expenses	true	"Expense card data"
//	@Success	201		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses [put]
func (h *expensesHandler) CreateExpenses(ctx *gin.Context) {
	var r req.Expenses
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.expensesService.InsertExpenses(ctx.Request.Context(), r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating expenses"})
		return
	}

	ctx.JSON(http.StatusCreated, res.DefaultResponse{
		Status:  string(res.Success),
		Message: "Expenses created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	})
}

// GetListMoneyCard godoc
//
//	@Summary	Get list of expense cards
//	@Tags		Expenses
//	@Produce	json
//	@Success	200	{object}	res.DefaultResponse
//	@Failure	500	{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses [get]
func (h *expensesHandler) GetListMoneyCard(ctx *gin.Context) {
	listData, err := h.expensesService.GetListMoneyCard(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching expenses"})
		return
	}
	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status: string(res.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData,
	})
}

// GetListMoneyCardDetail godoc
//
//	@Summary	Get expense detail list by card id
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.GetExpensesDetailById	true	"Card id"
//	@Success	200		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details [post]
func (h *expensesHandler) GetListMoneyCardDetail(ctx *gin.Context) {
	var r req.GetExpensesDetailById
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listData, err := h.expensesService.GetListMoneyCardDetail(ctx.Request.Context(), r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching expense details"})
		return
	}
	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status: string(res.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   listData,
	})
}

// UpsertExpensesDetail godoc
//
//	@Summary	Create or update an expense detail
//	@Description	Creates a new expense item when ID is omitted/zero, otherwise updates the existing item.
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.ExpensesDetail	true	"ExpensesDetail data"
//	@Success	200		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details [put]
func (h *expensesHandler) UpsertExpensesDetail(ctx *gin.Context) {
	var r req.ExpensesDetail
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID == 0 means create, otherwise update the existing detail.
	if r.ID == 0 {
		result, err := h.expensesService.InsertExpensesDetail(ctx.Request.Context(), r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating expenses detail"})
			return
		}
		ctx.JSON(http.StatusCreated, res.DefaultResponse{
			Status:  string(res.Success),
			Message: "Expenses detail created successfully",
			Date:    time.Now().Format("02/01/2006 15:04:05"),
			Data:    result,
		})
		return
	}

	result, err := h.expensesService.UpdateExpensesDetail(ctx.Request.Context(), r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating expenses detail"})
		return
	}
	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status:  string(res.Success),
		Message: "Expenses detail updated successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	})
}

// DeleteExpensesDetail godoc
//
//	@Summary	Delete an expense detail
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.DeleteExpensesDetailById	true	"ExpensesDetail id"
//	@Success	200		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details [delete]
func (h *expensesHandler) DeleteExpensesDetail(ctx *gin.Context) {
	var r req.DeleteExpensesDetailById
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, err := h.expensesService.DeleteExpensesDetail(ctx.Request.Context(), r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting expenses detail"})
		return
	}

	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status:  string(res.Success),
		Message: "Expenses detail deleted successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    ok,
	})
}
