package repository

import (
	"context"
	"log/slog"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/repository/interfaces"
	"gorm.io/gorm"
)

type expensesRepository struct {
	db *gorm.DB
}

func NewExpensesRepository(db *gorm.DB) interfaces.ExpensesRepository {
	return &expensesRepository{db: db}
}

func (r *expensesRepository) InsertExpenses(ctx context.Context, obj model.Expenses) (model.Expenses, error) {
	if err := r.db.WithContext(ctx).Create(&obj).Error; err != nil {
		slog.ErrorContext(ctx, "InsertExpenses: create failed", "error", err)
		return model.Expenses{}, err
	}
	return obj, nil
}

func (r *expensesRepository) InsertExpensesDetail(ctx context.Context, obj model.ExpensesDetail) (model.ExpensesDetail, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&obj).Error; err != nil {
			return err
		}
		var expenses model.Expenses
		if err := tx.First(&expenses, obj.ExpensesId).Error; err != nil {
			slog.ErrorContext(ctx, "InsertExpensesDetail: load parent expenses failed", "error", err)
			return err
		}
		newBalance := expenses.ExpensesBalance - obj.ExpensesAmount
		if err := tx.Model(&expenses).Updates(model.Expenses{ExpensesBalance: newBalance}).Error; err != nil {
			slog.ErrorContext(ctx, "InsertExpensesDetail: update balance failed", "error", err)
			return err
		}
		return nil
	})
	if err != nil {
		return model.ExpensesDetail{}, err
	}
	return obj, nil
}

func (r *expensesRepository) GetListMoneyCard(ctx context.Context) ([]model.Expenses, error) {
	rows, err := r.db.WithContext(ctx).Raw("SELECT * FROM expenses ORDER BY expenses_month ASC").Rows()
	if err != nil {
		slog.ErrorContext(ctx, "GetListMoneyCard: query failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var expensesData []model.Expenses
	for rows.Next() {
		var e model.Expenses
		if err := r.db.ScanRows(rows, &e); err != nil {
			slog.ErrorContext(ctx, "GetListMoneyCard: scan row failed", "error", err)
			return nil, err
		}
		expensesData = append(expensesData, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expensesData, nil
}

func (r *expensesRepository) GetListMoneyCardDetail(ctx context.Context, req req.GetExpensesDetailById) ([]model.ExpensesDetail, error) {
	rows, err := r.db.WithContext(ctx).Raw(
		"SELECT * FROM expenses_details WHERE expenses_id = ? ORDER BY created_at DESC",
		req.Id,
	).Rows()
	if err != nil {
		slog.ErrorContext(ctx, "GetListMoneyCardDetail: query failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []model.ExpensesDetail
	for rows.Next() {
		var d model.ExpensesDetail
		if err := r.db.ScanRows(rows, &d); err != nil {
			slog.ErrorContext(ctx, "GetListMoneyCardDetail: scan row failed", "error", err)
			return nil, err
		}
		details = append(details, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return details, nil
}

func (r *expensesRepository) UpdateExpensesDetail(ctx context.Context, req req.ExpensesDetail) (model.ExpensesDetail, error) {
	var detail model.ExpensesDetail
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&detail, req.ID).Error; err != nil {
			slog.ErrorContext(ctx, "UpdateExpensesDetail: load detail failed", "error", err)
			return err
		}
		var expenses model.Expenses
		if err := tx.First(&expenses, detail.ExpensesId).Error; err != nil {
			slog.ErrorContext(ctx, "UpdateExpensesDetail: load parent expenses failed", "error", err)
			return err
		}
		newBalance := expenses.ExpensesBalance + detail.ExpensesAmount - req.ExpensesAmount
		if err := tx.Model(&expenses).Updates(map[string]interface{}{"expenses_balance": newBalance}).Error; err != nil {
			slog.ErrorContext(ctx, "UpdateExpensesDetail: update balance failed", "error", err)
			return err
		}
		detail.ExpensesType = req.ExpensesType
		detail.ExpensesDesc = req.ExpensesDesc
		detail.ExpensesAmount = req.ExpensesAmount
		if err := tx.Save(&detail).Error; err != nil {
			slog.ErrorContext(ctx, "UpdateExpensesDetail: save detail failed", "error", err)
			return err
		}
		return nil
	})
	if err != nil {
		return model.ExpensesDetail{}, err
	}
	return detail, nil
}

func (r *expensesRepository) DeleteExpensesDetail(ctx context.Context, req req.DeleteExpensesDetailById) (bool, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var detail model.ExpensesDetail
		if err := tx.First(&detail, req.ID).Error; err != nil {
			slog.ErrorContext(ctx, "DeleteExpensesDetail: load detail failed", "error", err)
			return err
		}
		var expenses model.Expenses
		if err := tx.First(&expenses, detail.ExpensesId).Error; err != nil {
			slog.ErrorContext(ctx, "DeleteExpensesDetail: load parent expenses failed", "error", err)
			return err
		}
		newBalance := expenses.ExpensesBalance + detail.ExpensesAmount
		if err := tx.Model(&expenses).Updates(map[string]interface{}{"expenses_balance": newBalance}).Error; err != nil {
			slog.ErrorContext(ctx, "DeleteExpensesDetail: update balance failed", "error", err)
			return err
		}
		if err := tx.Unscoped().Delete(&detail).Error; err != nil {
			slog.ErrorContext(ctx, "DeleteExpensesDetail: delete failed", "error", err)
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *expensesRepository) CountExpenses(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Expenses{}).Count(&count).Error; err != nil {
		slog.ErrorContext(ctx, "CountExpenses: count failed", "error", err)
		return 0, err
	}
	return count, nil
}

func (r *expensesRepository) DeleteOldestExpenses(ctx context.Context) error {
	var expenses model.Expenses
	if err := r.db.WithContext(ctx).Order("id").First(&expenses).Error; err != nil {
		slog.ErrorContext(ctx, "DeleteOldestExpenses: find oldest failed", "error", err)
		return err
	}
	if err := r.db.WithContext(ctx).Unscoped().Delete(&expenses).Error; err != nil {
		slog.ErrorContext(ctx, "DeleteOldestExpenses: delete failed", "error", err)
		return err
	}
	return nil
}
