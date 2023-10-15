package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
)

type ExpensesServices interface {
	InsertExpenses(req req_dtos.Expenses) models.Expenses
	InsertExpensesDetail(req req_dtos.ExpensesDetail) models.ExpensesDetail
	GetListMoneyCard() res_dtos.ExpensesCard
	GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail
}
