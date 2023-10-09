package dtos

type Expenses struct {
	ExpensesMonth int32   `json:"expensesMonth"`
	ExpensesYear  int32   `json:"expensesYear"`
	ExpensesMoney float32 `json:"expensesMoney"`
}

type ExpensesDetail struct {
	ExpensesId     int32   `json:"expensesId"`
	ExpensesType   string  `json:"expensesType"`
	ExpensesDesc   string  `json:"expensesDesc"`
	ExpensesAmount float32 `json:"expensesAmount"`
}

type GetExpensesDetailById struct {
	Id int `json:"id"`
}
