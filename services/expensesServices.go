package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

type ExpensesServices interface {
	InsertExpenses(req req_dtos.Expenses) models.Expenses
	InsertExpensesDetail(c *gin.Context) models.ExpensesDetail
	GetListMoneyCard(c *gin.Context) res_dtos.ExpensesCard
	GetListMoneyCardDetail(c *gin.Context) []models.ExpensesDetail
}
