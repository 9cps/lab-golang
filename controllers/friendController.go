package controllers

import (
	"time"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
)

func CreateFriend(c *gin.Context) {
	var req req_dtos.Friend
	// Use c.ShouldBindJSON to bind the request body to the DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request, invalid JSON data",
		})
		return
	}

	obj := models.Friend{
		F_NAME: req.F_NAME,
		L_NAME: req.L_NAME,
		TEL_NO: req.TEL_NO,
	}

	result := initializers.DB.Create(&obj)

	if result.Error != nil {
		c.Status(400)
		return
	}

	response := res_dtos.DefaultResponse{
		Status:  "success",
		Message: "Friend created successfully",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	}

	c.JSON(200, response)
}
