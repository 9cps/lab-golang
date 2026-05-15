package interfaces

import (
	"context"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
)

type ExpensesService interface {
	InsertExpenses(ctx context.Context, r req.Expenses) (res.ExpensesResponse, error)
	InsertExpensesDetail(ctx context.Context, r req.ExpensesDetail) (res.ExpensesDetailResponse, error)
	GetListMoneyCard(ctx context.Context) (res.ExpensesCard, error)
	GetListMoneyCardDetail(ctx context.Context, r req.GetExpensesDetailById) ([]res.ExpensesDetailResponse, error)
	UpdateExpensesDetail(ctx context.Context, r req.UpdateExpensesDetail) (res.ExpensesDetailResponse, error)
	DeleteExpensesDetail(ctx context.Context, r req.DeleteExpensesDetailById) (bool, error)
}
