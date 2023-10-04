package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

func InsertExpenses(c *gin.Context) models.Expenses {
	var req req_dtos.Expenses

	// Map req to model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return models.Expenses{}
	}

	isCondition := DeleteExpensesIfCountExceeds(c, 6)

	if isCondition {
		obj := models.Expenses{
			ExpensesMonth:   req.ExpensesMonth,
			ExpensesYear:    req.ExpensesYear,
			ExpensesMoney:   req.ExpensesMoney, // จำนวนเงินตั้งต้น
			ExpensesBalance: req.ExpensesMoney, // จำนวนเงินคงเหลือ
		}

		result := initializers.DB.Create(&obj)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"error": "Failed to insert expenses data",
			})
			return models.Expenses{}
		}
		return obj
	}
	return models.Expenses{}
}

func InsertExpensesDetail(c *gin.Context) models.ExpensesDetail {
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

func GetListMoneyCard(c *gin.Context) []models.Expenses {
	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM expenses ORDER BY expenses_month ASC").Rows()

	// Check for errors
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to execute SQL query",
		})
		return []models.Expenses{} // Return an empty slice in case of error
	}
	defer rows.Close()

	var expenses models.Expenses
	var expensesData []models.Expenses

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
