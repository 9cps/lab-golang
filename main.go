package main

//import "fmt"
import (
	"github.com/9cps/api-go-gin/controllers"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConncetDatabse()
}

func main() {
	r := gin.Default()

	// HealtCheck
	r.GET("/HealthCheckAPI", controllers.HealthCheckAPI)
	r.GET("/HealthCheckDB", controllers.HealthCheckDB)

	// Expenses
	r.PUT("/CreateExpenses", controllers.CreateExpenses)
	r.PUT("/CreateExpensesDetail", controllers.CreateExpensesDetail)
	r.GET("/GetListMoneyCard", controllers.GetListMoneyCard)
	// r.DELETE("/DeleteFriend/:id", controllers.DeleteFriend)
	r.Run() // listen and serve on port .env
}
