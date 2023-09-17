package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB // Declare the package-level db variable

func ConncetDatabse() {
	LoadEnv()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("SERVER_NAME"),
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("SERVER_PORT"),
	)

	// Open a connection to the database using GORM and the SQL Server driver.
	var err error // Declare an err variable to store the error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set logger mode as needed.
	})

	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}

	fmt.Println("Connected to the database")
}
