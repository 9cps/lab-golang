package dtos

import "time"

// ExpensesCard is the response type for GetListMoneyCard.
type ExpensesCard struct {
	TotalBalance         float64    `json:"totalBalance"`
	PercentBalance       float64    `json:"percentBalance"`
	PercentSpendingMonth float64    `json:"percentSpendingMonth"`
	Data                 []Expenses `json:"data"`
}

// Expenses is used as an element inside ExpensesCard.
type Expenses struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	ExpensesMonth   int32     `json:"expensesMonth"`
	ExpensesYear    int32     `json:"expensesYear"`
	ExpensesMoney   float64   `json:"expensesMoney"`
	ExpensesBalance float64   `json:"expensesBalance"`
	TotalSpending   float64   `json:"totalSpending"`
}

// ExpensesResponse is the response DTO for a single Expenses record (create).
type ExpensesResponse struct {
	ID              uint    `json:"id"`
	ExpensesMonth   int32   `json:"expensesMonth"`
	ExpensesYear    int32   `json:"expensesYear"`
	ExpensesMoney   float64 `json:"expensesMoney"`
	ExpensesBalance float64 `json:"expensesBalance"`
}

// ExpensesDetailResponse is the response DTO for a single ExpensesDetail record.
type ExpensesDetailResponse struct {
	ID             uint    `json:"id"`
	ExpensesId     int32   `json:"expensesId"`
	ExpensesType   string  `json:"expensesType"`
	ExpensesDesc   string  `json:"expensesDesc"`
	ExpensesAmount float64 `json:"expensesAmount"`
}
