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
	r.POST("/CreateFriend", controllers.CreateFriend)
	r.Run() // listen and serve on port .env
}
