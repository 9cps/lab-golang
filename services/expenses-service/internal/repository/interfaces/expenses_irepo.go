package interfaces

import (
	"context"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
)

type ExpensesRepository interface {
	InsertExpenses(ctx context.Context, obj model.Expenses) (model.Expenses, error)
	InsertExpensesDetail(ctx context.Context, obj model.ExpensesDetail) (model.ExpensesDetail, error)
	GetListMoneyCard(ctx context.Context) ([]model.Expenses, error)
	GetListMoneyCardDetail(ctx context.Context, r req.GetExpensesDetailById) ([]model.ExpensesDetail, error)
	UpdateExpensesDetail(ctx context.Context, r req.ExpensesDetail) (model.ExpensesDetail, error)
	DeleteExpensesDetail(ctx context.Context, r req.DeleteExpensesDetailById) (bool, error)
	CountExpenses(ctx context.Context) (int64, error)
	DeleteOldestExpenses(ctx context.Context) error
}
