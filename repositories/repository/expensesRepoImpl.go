package repositories

import (
	"fmt"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/9cps/api-go-gin/repositories/interfaces"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpensesRepositoryImpl struct {
	Db *gorm.DB
}

func NewExpensesRepositoryImpl(Db *gorm.DB) interfaces.IExpensesRepository {
	return &ExpensesRepositoryImpl{Db: Db}
}

func (r *ExpensesRepositoryImpl) InsertExpenses(req models.Expenses) models.Expenses {
	result := initializers.DB.Create(&req)
	if result.Error != nil {
		return models.Expenses{}
	}
	return req
}

func (r *ExpensesRepositoryImpl) InsertExpensesDetail(req models.ExpensesDetail) models.ExpensesDetail {
	// Create the ExpensesDetail record
	result := initializers.DB.Create(&req)
	if result.Error != nil {
		return models.ExpensesDetail{}
	}

	// Update the Expenses record's ExpensesBalance
	var expenses models.Expenses
	if err := initializers.DB.First(&expenses, req.ExpensesId).Error; err != nil {
		fmt.Printf("Error creating ExpensesDetail: %v", result.Error)
		return models.ExpensesDetail{}
	}

	calBalance := expenses.ExpensesBalance - req.ExpensesAmount

	if err := initializers.DB.Model(&expenses).Updates(models.Expenses{ExpensesBalance: calBalance}).Error; err != nil {
		fmt.Printf("Error update ExpensesDetail: %v", result.Error)
		return models.ExpensesDetail{}
	}

	return req
}

func DeleteExpensesIfCountExceeds(c *gin.Context, threshold int64) bool {
	// Count the number of records
	var count int64
	if err := initializers.DB.Model(&models.Expenses{}).Count(&count).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to count expenses records",
		})
		return false
	}

	if count >= threshold {
		// Find and delete the record with the smallest ID
		var expenses models.Expenses
		if err := initializers.DB.Order("id").First(&expenses).Error; err != nil {
			c.JSON(400, gin.H{
				"error": "Failed to find the record to delete",
			})
			return false
		}

		// ใช้ลบ Unscoped เนื่องจากพื้น db มีจำกัด
		if err := initializers.DB.Unscoped().Delete(&expenses).Error; err != nil {
			c.JSON(400, gin.H{
				"error": "Failed to delete the record",
			})
			return false
		}
	}
	return true
}

func (r *ExpensesRepositoryImpl) GetListMoneyCard() res_dtos.ExpensesCard {
	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM expenses ORDER BY expenses_month ASC").Rows()

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
		initializers.DB.ScanRows(rows, &expenses)
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

func (r *ExpensesRepositoryImpl) GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail {
	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM expenses_details WHERE expenses_id = ? ORDER BY created_at DESC", req.Id).Rows()

	// Check for errors
	if err != nil {
		return []models.ExpensesDetail{} // Return an empty slice in case of error
	}
	defer rows.Close()

	var expenses models.ExpensesDetail
	var expensesData []models.ExpensesDetail

	for rows.Next() {
		// Scan the result into the friend struct
		initializers.DB.ScanRows(rows, &expenses)
		expensesData = append(expensesData, expenses)
	}

	return expensesData
}
