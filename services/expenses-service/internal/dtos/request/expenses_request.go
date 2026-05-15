package dtos

type Expenses struct {
	ExpensesMonth int32   `json:"expensesMonth" binding:"required"`
	ExpensesYear  int32   `json:"expensesYear"  binding:"required"`
	ExpensesMoney float64 `json:"expensesMoney"  binding:"required"`
}

type ExpensesDetail struct {
	ExpensesId     int32   `json:"expensesId"     binding:"required"`
	ExpensesType   string  `json:"expensesType"   binding:"required"`
	ExpensesDesc   string  `json:"expensesDesc"`
	ExpensesAmount float64 `json:"expensesAmount" binding:"required"`
}

type GetExpensesDetailById struct {
	Id int `form:"id" json:"id" binding:"required"`
}

type UpdateExpensesDetail struct {
	Id             uint    `json:"id"`
	ExpensesType   string  `json:"expensesType"   binding:"required"`
	ExpensesDesc   string  `json:"expensesDesc"`
	ExpensesAmount float64 `json:"expensesAmount" binding:"required"`
}

type DeleteExpensesDetailById struct {
	Id uint `json:"id"`
}
