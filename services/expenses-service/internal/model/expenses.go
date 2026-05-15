package model

import "gorm.io/gorm"

type ExpensesDetail struct {
	gorm.Model
	ExpensesId     int32  `gorm:"not null"`
	ExpensesType   string `gorm:"not null"`
	ExpensesDesc   string
	ExpensesAmount float64 `gorm:"not null"`
}

type Expenses struct {
	gorm.Model
	ExpensesMonth   int32   `gorm:"not null"`
	ExpensesYear    int32   `gorm:"not null"`
	ExpensesMoney   float64 `gorm:"not null"`
	ExpensesBalance float64 `gorm:"not null"`
}
