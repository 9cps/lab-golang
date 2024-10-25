package repositories

import (
	"fmt"
	"os"

	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/repositories/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	initializers.LoadEnv()
}

type HealthCheckRepositoryImpl struct {
	Db *gorm.DB
}

func NewHealthCheckRepositoryImpl(Db *gorm.DB) interfaces.IHealthCheckRepository {
	return &HealthCheckRepositoryImpl{Db: Db}
}

func (r *HealthCheckRepositoryImpl) HealthCheckDB() bool {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("SERVER_NAME"),
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("SERVER_PORT"),
	)

	// Open a connection to the database using GORM and the SQL Server driver.
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set logger mode as needed.
	})

	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return false // Return false to indicate a connection failure
	}

	return true // Return true to indicate a successful connection
}
