package repository

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
)

type ExpensesRopository interface {
	InsertExpenses(req models.Expenses) models.Expenses
	InsertExpensesDetail(req models.ExpensesDetail) models.ExpensesDetail
	GetListMoneyCard() res_dtos.ExpensesCard
	GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail
}
