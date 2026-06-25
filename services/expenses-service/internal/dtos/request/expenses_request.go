package dtos

// Expenses is the request body for creating a new monthly expense card.
// Route: PUT /expenses
type Expenses struct {
	ExpensesMonth int32   `json:"ExpensesMonth" binding:"required"`
	ExpensesYear  int32   `json:"ExpensesYear"  binding:"required"`
	ExpensesMoney float64 `json:"ExpensesMoney" binding:"required"`
}

// ExpensesDetail is the request body for creating OR updating an expense item.
// Route: PUT /expenses/details
//
// When ID is zero the item is created; when ID is non-zero the existing item
// is updated. ExpensesId identifies the parent card.
type ExpensesDetail struct {
	ID             uint    `json:"ID"`
	ExpensesId     int32   `json:"ExpensesId"     binding:"required"`
	ExpensesType   string  `json:"ExpensesType"   binding:"required"`
	ExpensesDesc   string  `json:"ExpensesDesc"`
	ExpensesAmount float64 `json:"ExpensesAmount" binding:"required"`
}

// GetExpensesDetailById is the request body for listing the items of a card.
// Route: POST /expenses/details
type GetExpensesDetailById struct {
	Id int `json:"id" binding:"required"`
}

// DeleteExpensesDetailById is the request body for deleting an expense item.
// Route: DELETE /expenses/details
type DeleteExpensesDetailById struct {
	ID         uint  `json:"ID" binding:"required"`
	ExpensesId int32 `json:"ExpensesId"`
}
