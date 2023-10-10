package dtos

import "gorm.io/gorm"

type ExpensesCard struct {
	TotalBalance         float32
	PercentBalance       float32
	PercentSpendingMonth float32
	Data                 []Expenses
}

type Expenses struct {
	gorm.Model
	ExpensesMonth   int32
	ExpensesYear    int32
	ExpensesMoney   float32
	ExpensesBalance float32
	TotalSpending   float32
}
