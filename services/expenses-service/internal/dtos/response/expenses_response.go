package dtos

import "time"

// ExpensesCard is the response payload for GetListMoneyCard.
// It is nested under the top-level `data` field of DefaultResponse so the UI
// reads it as `data.Data` and `data.TotalBalance`.
type ExpensesCard struct {
	Data         []Expenses `json:"Data"`
	TotalBalance float64    `json:"TotalBalance"`
}

// Expenses is a single monthly expense card element inside ExpensesCard.
type Expenses struct {
	ID              uint    `json:"ID"`
	ExpensesMonth   int32   `json:"ExpensesMonth"`
	ExpensesYear    int32   `json:"ExpensesYear"`
	ExpensesBalance float64 `json:"ExpensesBalance"`
	ExpensesMoney   float64 `json:"ExpensesMoney"`
	TotalSpending   float64 `json:"TotalSpending"`
}

// ExpensesResponse is the response DTO for a single Expenses record (create).
type ExpensesResponse struct {
	ID              uint    `json:"ID"`
	ExpensesMonth   int32   `json:"ExpensesMonth"`
	ExpensesYear    int32   `json:"ExpensesYear"`
	ExpensesMoney   float64 `json:"ExpensesMoney"`
	ExpensesBalance float64 `json:"ExpensesBalance"`
}

// ExpensesDetailResponse is the response DTO for a single ExpensesDetail record.
type ExpensesDetailResponse struct {
	ID             uint      `json:"ID"`
	ExpensesId     int32     `json:"ExpensesId"`
	ExpensesType   string    `json:"ExpensesType"`
	ExpensesDesc   string    `json:"ExpensesDesc"`
	ExpensesAmount float64   `json:"ExpensesAmount"`
	CreatedAt      time.Time `json:"CreatedAt"`
}
