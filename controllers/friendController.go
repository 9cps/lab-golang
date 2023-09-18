package controllers

import (
	"time"

	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	services "github.com/9cps/api-go-gin/services"
	"github.com/gin-gonic/gin"
)

func CreateFriend(c *gin.Context) {
	result := services.InsertFriend(c)
	if result == (models.Friend{}) {
		c.JSON(400, gin.H{
			"error": "Error creating friend",
		})
		return
	}
	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "Friend created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
		Data:    result,
	}

	c.JSON(200, response)
}

func GetFriend(c *gin.Context) {
	// Find friend by keyword
	friendData := services.FindFriend(c)
	response := res_dtos.DefaultResponse{
		Status: string(res_dtos.Success),
		Date:   time.Now().Format("02/01/2006 15:04:05"),
		Data:   friendData, // Data retrieved from the SQL query
	}

	// Return data
	c.JSON(200, response)
}
