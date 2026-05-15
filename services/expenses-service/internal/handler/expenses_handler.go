package handler

import (
	"net/http"
	"strconv"
	"time"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type ExpensesHandler interface {
	CreateExpenses(ctx *gin.Context)
	CreateExpensesDetail(ctx *gin.Context)
	GetListMoneyCard(ctx *gin.Context)
	GetListMoneyCardDetail(ctx *gin.Context)
	UpdateExpensesDetail(ctx *gin.Context)
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
//	@Summary	Create expenses
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.Expenses	true	"Expense data"
//	@Success	201		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses [post]
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

// CreateExpensesDetail godoc
//
//	@Summary	Create expenses detail
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		body	body		req.ExpensesDetail	true	"ExpensesDetail data"
//	@Success	201		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details [post]
func (h *expensesHandler) CreateExpensesDetail(ctx *gin.Context) {
	var r req.ExpensesDetail
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
//	@Summary	Get expense detail list by expense id
//	@Tags		Expenses
//	@Produce	json
//	@Param		id	query		int	true	"Expense ID"
//	@Success	200	{object}	res.DefaultResponse
//	@Failure	400	{object}	gin.H
//	@Failure	500	{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details [get]
func (h *expensesHandler) GetListMoneyCardDetail(ctx *gin.Context) {
	var r req.GetExpensesDetailById
	if err := ctx.ShouldBindQuery(&r); err != nil {
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

// UpdateExpensesDetail godoc
//
//	@Summary	Update expense detail
//	@Tags		Expenses
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"ExpensesDetail ID"
//	@Param		body	body		req.UpdateExpensesDetail	true	"Update data"
//	@Success	200		{object}	res.DefaultResponse
//	@Failure	400		{object}	gin.H
//	@Failure	500		{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details/{id} [put]
func (h *expensesHandler) UpdateExpensesDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var r req.UpdateExpensesDetail
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.Id = uint(id)

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
//	@Summary	Delete expense detail
//	@Tags		Expenses
//	@Produce	json
//	@Param		id	path		int	true	"ExpensesDetail ID"
//	@Success	200	{object}	res.DefaultResponse
//	@Failure	400	{object}	gin.H
//	@Failure	500	{object}	gin.H
//	@Security	BearerAuth
//	@Router		/expenses/details/{id} [delete]
func (h *expensesHandler) DeleteExpensesDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ok, err := h.expensesService.DeleteExpensesDetail(ctx.Request.Context(), req.DeleteExpensesDetailById{Id: uint(id)})
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
