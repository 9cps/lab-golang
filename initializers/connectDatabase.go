package initializers

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the shared *gorm.DB set by ConnectDatabase. Kept exported for
// backwards compatibility with existing call sites; prefer passing the *gorm.DB
// returned from ConnectDatabase explicitly.
var DB *gorm.DB

// ConnectDatabase opens a PostgreSQL connection using environment variables
// and configures the connection pool. It sets the package-level DB variable
// and also returns the *gorm.DB so callers can wire dependencies explicitly.
func ConnectDatabase() (*gorm.DB, error) {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("SERVER_NAME"),
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("SERVER_PORT"),
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	DB = db
	fmt.Println("Connected to the database")
	return db, nil
}

// ConncetDatabse is a deprecated alias for the misspelled original name.
// New code should use ConnectDatabase instead.
//
// Deprecated: use ConnectDatabase instead.
func ConncetDatabse() {
	if _, err := ConnectDatabase(); err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
	}
}
