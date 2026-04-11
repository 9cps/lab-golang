package repositories

import (
	"fmt"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	"github.com/9cps/api-go-gin/repositories/interfaces"
	"gorm.io/gorm"
)

type ExpensesRepositoryImpl struct {
	Db *gorm.DB
}

func NewExpensesRepositoryImpl(Db *gorm.DB) interfaces.IExpensesRepository {
	return &ExpensesRepositoryImpl{Db: Db}
}

func (r *ExpensesRepositoryImpl) InsertExpenses(req models.Expenses) models.Expenses {
	result := r.Db.Create(&req)
	if result.Error != nil {
		return models.Expenses{}
	}
	return req
}

func (r *ExpensesRepositoryImpl) InsertExpensesDetail(req models.ExpensesDetail) models.ExpensesDetail {
	// Create the ExpensesDetail record
	result := r.Db.Create(&req)
	if result.Error != nil {
		return models.ExpensesDetail{}
	}

	// Update the Expenses record's ExpensesBalance
	var expenses models.Expenses
	if err := r.Db.First(&expenses, req.ExpensesId).Error; err != nil {
		fmt.Printf("Error loading parent Expenses: %v\n", err)
		return models.ExpensesDetail{}
	}

	calBalance := expenses.ExpensesBalance - req.ExpensesAmount

	if err := r.Db.Model(&expenses).Updates(models.Expenses{ExpensesBalance: calBalance}).Error; err != nil {
		fmt.Printf("Error updating Expenses balance: %v\n", err)
		return models.ExpensesDetail{}
	}

	return req
}

func (r *ExpensesRepositoryImpl) GetListMoneyCard() res_dtos.ExpensesCard {
	// SQL Query
	rows, err := r.Db.Raw("SELECT * FROM expenses ORDER BY expenses_month ASC").Rows()

	// Check for errors
	if err != nil {
		fmt.Printf("Error sql GetListMoneyCard: %v", err)
		return res_dtos.ExpensesCard{} // Return an empty slice in case of error
	}
	defer rows.Close()

	var expensesData []res_dtos.Expenses

	sumBalance := float32(0)
	for rows.Next() {
		var expenses res_dtos.Expenses
		// Scan the result into the friend struct
		r.Db.ScanRows(rows, &expenses)
		sumBalance += expenses.ExpensesBalance
		// Calculate spending for each row
		spending := expenses.ExpensesMoney - expenses.ExpensesBalance
		expenses.TotalSpending += spending
		expensesData = append(expensesData, expenses)
	}

	var expensesCard res_dtos.ExpensesCard
	expensesCard.Data = expensesData
	expensesCard.TotalBalance = sumBalance

	return expensesCard
}

func (r *ExpensesRepositoryImpl) UpdateExpensesDetail(req req_dtos.UpdateExpensesDetail) models.ExpensesDetail {
	var detail models.ExpensesDetail
	if err := r.Db.First(&detail, req.Id).Error; err != nil {
		fmt.Printf("Error loading ExpensesDetail: %v\n", err)
		return models.ExpensesDetail{}
	}

	var expenses models.Expenses
	if err := r.Db.First(&expenses, detail.ExpensesId).Error; err != nil {
		fmt.Printf("Error loading parent Expenses: %v\n", err)
		return models.ExpensesDetail{}
	}

	newBalance := expenses.ExpensesBalance + detail.ExpensesAmount - req.ExpensesAmount
	if err := r.Db.Model(&expenses).Updates(map[string]interface{}{"expenses_balance": newBalance}).Error; err != nil {
		fmt.Printf("Error updating Expenses balance: %v\n", err)
		return models.ExpensesDetail{}
	}

	detail.ExpensesType = req.ExpensesType
	detail.ExpensesDesc = req.ExpensesDesc
	detail.ExpensesAmount = req.ExpensesAmount
	if err := r.Db.Save(&detail).Error; err != nil {
		fmt.Printf("Error saving ExpensesDetail: %v\n", err)
		return models.ExpensesDetail{}
	}
	return detail
}

func (r *ExpensesRepositoryImpl) DeleteExpensesDetail(req req_dtos.DeleteExpensesDetailById) bool {
	var detail models.ExpensesDetail
	if err := r.Db.First(&detail, req.Id).Error; err != nil {
		fmt.Printf("Error loading ExpensesDetail: %v\n", err)
		return false
	}

	var expenses models.Expenses
	if err := r.Db.First(&expenses, detail.ExpensesId).Error; err != nil {
		fmt.Printf("Error loading parent Expenses: %v\n", err)
		return false
	}

	newBalance := expenses.ExpensesBalance + detail.ExpensesAmount
	if err := r.Db.Model(&expenses).Updates(map[string]interface{}{"expenses_balance": newBalance}).Error; err != nil {
		fmt.Printf("Error updating Expenses balance: %v\n", err)
		return false
	}

	if err := r.Db.Unscoped().Delete(&detail).Error; err != nil {
		fmt.Printf("Error deleting ExpensesDetail: %v\n", err)
		return false
	}
	return true
}

func (r *ExpensesRepositoryImpl) GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail {
	// SQL Query
	rows, err := r.Db.Raw("SELECT * FROM expenses_details WHERE expenses_id = ? ORDER BY created_at DESC", req.Id).Rows()

	// Check for errors
	if err != nil {
		return []models.ExpensesDetail{} // Return an empty slice in case of error
	}
	defer rows.Close()

	var expenses models.ExpensesDetail
	var expensesData []models.ExpensesDetail

	for rows.Next() {
		// Scan the result into the friend struct
		r.Db.ScanRows(rows, &expenses)
		expensesData = append(expensesData, expenses)
	}

	return expensesData
}
