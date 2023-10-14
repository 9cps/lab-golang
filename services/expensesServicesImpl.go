package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/9cps/api-go-gin/repository"
	"github.com/gin-gonic/gin"
)

type ExpensesServiceImpl struct {
	ExpensesRopository repository.ExpensesRopository
}

func NewExpensesServiceImpl(expensesRopository repository.ExpensesRopository) ExpensesServices {
	return &ExpensesServiceImpl{
		ExpensesRopository: expensesRopository,
	}
}

func (s *ExpensesServiceImpl) InsertExpenses(req req_dtos.Expenses) models.Expenses {

	isCondition := DeleteExpensesIfCountExceeds(req, 6)

	if isCondition {
		obj := models.Expenses{
			ExpensesMonth:   req.ExpensesMonth,
			ExpensesYear:    req.ExpensesYear,
			ExpensesMoney:   req.ExpensesMoney, // จำนวนเงินตั้งต้น
			ExpensesBalance: req.ExpensesMoney, // จำนวนเงินคงเหลือ
		}

		result := initializers.DB.Create(&obj)
		if result.Error != nil {
			return models.Expenses{}
		}
		return obj
	}
	return models.Expenses{}
}

func (s *ExpensesServiceImpl) InsertExpensesDetail(c *gin.Context) models.ExpensesDetail {
	var req req_dtos.ExpensesDetail

	// Map req to model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return models.ExpensesDetail{}
	}

	obj := models.ExpensesDetail{
		ExpensesId:     req.ExpensesId,
		ExpensesType:   req.ExpensesType,
		ExpensesDesc:   req.ExpensesDesc,
		ExpensesAmount: req.ExpensesAmount,
	}

	// Create the ExpensesDetail record
	result := initializers.DB.Create(&obj)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "Failed to insert expenses detail data",
		})
		return models.ExpensesDetail{}
	}

	// Update the Expenses record's ExpensesBalance
	var expenses models.Expenses
	if err := initializers.DB.First(&expenses, req.ExpensesId).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Expenses record not found",
		})
		return models.ExpensesDetail{}
	}

	calBalance := expenses.ExpensesBalance - req.ExpensesAmount

	if err := initializers.DB.Model(&expenses).Updates(models.Expenses{ExpensesBalance: calBalance}).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to update Expenses record",
		})
		return models.ExpensesDetail{}
	}

	return obj
}

func DeleteExpensesIfCountExceeds(req req_dtos.Expenses, threshold int64) bool {
	// Count the number of records
	var count int64
	if err := initializers.DB.Model(&models.Expenses{}).Count(&count).Error; err != nil {
		return false
	}

	if count >= threshold {
		// Find and delete the record with the smallest ID
		var expenses models.Expenses
		if err := initializers.DB.Order("id").First(&expenses).Error; err != nil {
			return false
		}

		// ใช้ลบ Unscoped เนื่องจากพื้น db มีจำกัด
		if err := initializers.DB.Unscoped().Delete(&expenses).Error; err != nil {
			return false
		}
	}
	return true
}

func (s *ExpensesServiceImpl) GetListMoneyCard(c *gin.Context) res_dtos.ExpensesCard {
	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM expenses ORDER BY expenses_month ASC").Rows()

	// Check for errors
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to execute SQL query",
		})
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

func (s *ExpensesServiceImpl) GetListMoneyCardDetail(c *gin.Context) []models.ExpensesDetail {
	var req req_dtos.GetExpensesDetailById
	// Map req to model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return []models.ExpensesDetail{} // Return an empty slice in case of error
	}

	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM expenses_details WHERE expenses_id = ? ORDER BY created_at DESC", req.Id).Rows()

	// Check for errors
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to execute SQL query",
		})
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

// func FindFriend(c *gin.Context) []models.Friend {
// 	var req req_dtos.GetFriend
// 	// Map req to model
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Bad request, invalid JSON data",
// 		})
// 		return []models.Friend{} // Return an empty slice in case of error
// 	}

// 	// SQL Query
// 	rows, err := initializers.DB.Raw("SELECT * FROM friends WHERE F_NAME LIKE ? OR L_NAME LIKE ?", "%"+req.KEYWORD+"%", "%"+req.KEYWORD+"%").Rows()

// 	// Check for errors
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Failed to execute SQL query",
// 		})
// 		return []models.Friend{} // Return an empty slice in case of error
// 	}
// 	defer rows.Close()

// 	var friend models.Friend // Assuming FriendModel is the struct for your friend data
// 	var friendData []models.Friend

// 	for rows.Next() {
// 		// Scan the result into the friend struct
// 		initializers.DB.ScanRows(rows, &friend)
// 		friendData = append(friendData, friend)
// 	}

// 	return friendData
// }

// func DeleteFriend(c *gin.Context) bool {
// 	id := c.Param("id")

// 	// Delete by id
// 	initializers.DB.Delete(&models.Friend{}, id)

// 	return true
// }
