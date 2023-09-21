package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

func InsertFriend(c *gin.Context) models.Friend {
	var req req_dtos.Friend
	// Map req to model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return models.Friend{}
	}

	obj := models.Friend{
		F_NAME: req.F_NAME,
		L_NAME: req.L_NAME,
		TEL_NO: req.TEL_NO,
	}

	result := initializers.DB.Create(&obj)
	if result.Error != nil {
		c.Status(400)
		return models.Friend{}
	}

	return obj
}

func FindFriend(c *gin.Context) []models.Friend {
	var req req_dtos.GetFriend
	// Map req to model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return []models.Friend{} // Return an empty slice in case of error
	}

	// SQL Query
	rows, err := initializers.DB.Raw("SELECT * FROM friends WHERE F_NAME LIKE ? OR L_NAME LIKE ?", "%"+req.KEYWORD+"%", "%"+req.KEYWORD+"%").Rows()

	// Check for errors
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to execute SQL query",
		})
		return []models.Friend{} // Return an empty slice in case of error
	}
	defer rows.Close()

	var friend models.Friend // Assuming FriendModel is the struct for your friend data
	var friendData []models.Friend

	for rows.Next() {
		// Scan the result into the friend struct
		initializers.DB.ScanRows(rows, &friend)
		friendData = append(friendData, friend)
	}

	return friendData
}

func DeleteFriend(c *gin.Context) bool {
	id := c.Param("id")

	// Delete by id
	initializers.DB.Delete(&models.Friend{}, id)

	return true
}
