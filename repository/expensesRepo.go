package repository

import (
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

type ExpensesRopository interface {
	InsertExpenses(c *gin.Context) models.Expenses
	InsertExpensesDetail(c *gin.Context) models.ExpensesDetail
	GetListMoneyCard(c *gin.Context) res_dtos.ExpensesCard
	GetListMoneyCardDetail(c *gin.Context) []models.ExpensesDetail
}
