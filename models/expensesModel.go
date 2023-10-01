package models

import "gorm.io/gorm"

type ExpensesDetail struct {
	gorm.Model
	ExpensesId     int32
	ExpensesType   string
	ExpensesDesc   string
	ExpensesAmount float32
}

type Expenses struct {
	gorm.Model
	ExpensesMonth   int32
	ExpensesYear    int32
	ExpensesMoney   float32
	ExpensesBalance float32
}
