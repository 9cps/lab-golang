package service

import (
	"context"
	"fmt"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
	repoInterfaces "github.com/9cps/api-go-gin/services/expenses-service/internal/repository/interfaces"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service/interfaces"
)

type expensesService struct {
	expensesRepo repoInterfaces.ExpensesRepository
}

func NewExpensesService(expensesRepo repoInterfaces.ExpensesRepository) interfaces.ExpensesService {
	return &expensesService{expensesRepo: expensesRepo}
}

func (s *expensesService) deleteIfCountExceeds(ctx context.Context, threshold int64) error {
	count, err := s.expensesRepo.CountExpenses(ctx)
	if err != nil {
		return fmt.Errorf("count expenses: %w", err)
	}
	if count >= threshold {
		if err := s.expensesRepo.DeleteOldestExpenses(ctx); err != nil {
			return fmt.Errorf("delete oldest expenses: %w", err)
		}
	}
	return nil
}

func (s *expensesService) InsertExpenses(ctx context.Context, r req.Expenses) (res.ExpensesResponse, error) {
	if err := s.deleteIfCountExceeds(ctx, 6); err != nil {
		return res.ExpensesResponse{}, err
	}
	obj := model.Expenses{
		ExpensesMonth:   r.ExpensesMonth,
		ExpensesYear:    r.ExpensesYear,
		ExpensesMoney:   r.ExpensesMoney,
		ExpensesBalance: r.ExpensesMoney,
	}
	created, err := s.expensesRepo.InsertExpenses(ctx, obj)
	if err != nil {
		return res.ExpensesResponse{}, fmt.Errorf("insert expenses: %w", err)
	}
	return res.ExpensesResponse{
		ID:              created.ID,
		ExpensesMonth:   created.ExpensesMonth,
		ExpensesYear:    created.ExpensesYear,
		ExpensesMoney:   created.ExpensesMoney,
		ExpensesBalance: created.ExpensesBalance,
	}, nil
}

func (s *expensesService) InsertExpensesDetail(ctx context.Context, r req.ExpensesDetail) (res.ExpensesDetailResponse, error) {
	obj := model.ExpensesDetail{
		ExpensesId:     r.ExpensesId,
		ExpensesType:   r.ExpensesType,
		ExpensesDesc:   r.ExpensesDesc,
		ExpensesAmount: r.ExpensesAmount,
	}
	created, err := s.expensesRepo.InsertExpensesDetail(ctx, obj)
	if err != nil {
		return res.ExpensesDetailResponse{}, fmt.Errorf("insert expenses detail: %w", err)
	}
	return res.ExpensesDetailResponse{
		ID:             created.ID,
		ExpensesId:     created.ExpensesId,
		ExpensesType:   created.ExpensesType,
		ExpensesDesc:   created.ExpensesDesc,
		ExpensesAmount: created.ExpensesAmount,
	}, nil
}

func (s *expensesService) GetListMoneyCard(ctx context.Context) (res.ExpensesCard, error) {
	rows, err := s.expensesRepo.GetListMoneyCard(ctx)
	if err != nil {
		return res.ExpensesCard{}, fmt.Errorf("get list money card: %w", err)
	}

	var totalBalance float64
	data := make([]res.Expenses, 0, len(rows))
	for _, e := range rows {
		totalBalance += e.ExpensesBalance
		data = append(data, res.Expenses{
			ID:              e.ID,
			CreatedAt:       e.CreatedAt,
			ExpensesMonth:   e.ExpensesMonth,
			ExpensesYear:    e.ExpensesYear,
			ExpensesMoney:   e.ExpensesMoney,
			ExpensesBalance: e.ExpensesBalance,
			TotalSpending:   e.ExpensesMoney - e.ExpensesBalance,
		})
	}
	return res.ExpensesCard{TotalBalance: totalBalance, Data: data}, nil
}

func (s *expensesService) GetListMoneyCardDetail(ctx context.Context, r req.GetExpensesDetailById) ([]res.ExpensesDetailResponse, error) {
	rows, err := s.expensesRepo.GetListMoneyCardDetail(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("get list money card detail: %w", err)
	}
	result := make([]res.ExpensesDetailResponse, 0, len(rows))
	for _, d := range rows {
		result = append(result, res.ExpensesDetailResponse{
			ID:             d.ID,
			ExpensesId:     d.ExpensesId,
			ExpensesType:   d.ExpensesType,
			ExpensesDesc:   d.ExpensesDesc,
			ExpensesAmount: d.ExpensesAmount,
		})
	}
	return result, nil
}

func (s *expensesService) UpdateExpensesDetail(ctx context.Context, r req.UpdateExpensesDetail) (res.ExpensesDetailResponse, error) {
	updated, err := s.expensesRepo.UpdateExpensesDetail(ctx, r)
	if err != nil {
		return res.ExpensesDetailResponse{}, fmt.Errorf("update expenses detail: %w", err)
	}
	return res.ExpensesDetailResponse{
		ID:             updated.ID,
		ExpensesId:     updated.ExpensesId,
		ExpensesType:   updated.ExpensesType,
		ExpensesDesc:   updated.ExpensesDesc,
		ExpensesAmount: updated.ExpensesAmount,
	}, nil
}

func (s *expensesService) DeleteExpensesDetail(ctx context.Context, r req.DeleteExpensesDetailById) (bool, error) {
	ok, err := s.expensesRepo.DeleteExpensesDetail(ctx, r)
	if err != nil {
		return false, fmt.Errorf("delete expenses detail: %w", err)
	}
	return ok, nil
}
