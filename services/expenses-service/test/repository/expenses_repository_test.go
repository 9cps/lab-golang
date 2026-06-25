package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/repository"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testCtx = context.Background()

func newRepo(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	sqlDB, mk, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	require.NoError(t, err)
	return gdb, mk, func() { _ = sqlDB.Close() }
}

// ---- InsertExpenses ----

func TestRepo_InsertExpenses_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mk.ExpectCommit()

	r := repository.NewExpensesRepository(db)
	got, err := r.InsertExpenses(testCtx, model.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000, ExpensesBalance: 1000})

	assert.NoError(t, err)
	assert.Equal(t, float64(1000), got.ExpensesMoney)
	assert.NoError(t, mk.ExpectationsWereMet())
}

func TestRepo_InsertExpenses_Error(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses"`).WillReturnError(sql.ErrConnDone)
	mk.ExpectRollback()

	r := repository.NewExpensesRepository(db)
	got, err := r.InsertExpenses(testCtx, model.Expenses{ExpensesMoney: 1000})

	assert.Error(t, err)
	assert.Equal(t, model.Expenses{}, got)
}

// ---- InsertExpensesDetail ----

func TestRepo_InsertExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses_details"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 1000))
	mk.ExpectExec(`UPDATE "expenses" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repository.NewExpensesRepository(db)
	got, err := r.InsertExpensesDetail(testCtx, model.ExpensesDetail{ExpensesId: 10, ExpensesAmount: 200, ExpensesType: "food"})

	assert.NoError(t, err)
	assert.Equal(t, int32(10), got.ExpensesId)
	assert.NoError(t, mk.ExpectationsWereMet())
}

func TestRepo_InsertExpensesDetail_CreateError(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses_details"`).WillReturnError(sql.ErrConnDone)
	mk.ExpectRollback()

	r := repository.NewExpensesRepository(db)
	got, err := r.InsertExpensesDetail(testCtx, model.ExpensesDetail{ExpensesId: 10})

	assert.Error(t, err)
	assert.Equal(t, model.ExpensesDetail{}, got)
}

// ---- GetListMoneyCard ----

func TestRepo_GetListMoneyCard_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "expenses_month", "expenses_year", "expenses_money", "expenses_balance"}).
		AddRow(1, 1, 2026, 1000, 800).
		AddRow(2, 2, 2026, 2000, 1500)
	mk.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM expenses ORDER BY expenses_month ASC`)).WillReturnRows(rows)

	r := repository.NewExpensesRepository(db)
	got, err := r.GetListMoneyCard(testCtx)

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, float64(800), got[0].ExpensesBalance)
}

// ---- GetListMoneyCardDetail ----

func TestRepo_GetListMoneyCardDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "expenses_id", "expenses_type", "expenses_desc", "expenses_amount"}).
		AddRow(1, 5, "food", "lunch", 120).
		AddRow(2, 5, "bill", "water", 300)

	mk.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM expenses_details WHERE expenses_id = $1 ORDER BY created_at DESC`)).
		WithArgs(5).
		WillReturnRows(rows)

	r := repository.NewExpensesRepository(db)
	got, err := r.GetListMoneyCardDetail(testCtx, req.GetExpensesDetailById{Id: 5})

	assert.NoError(t, err)
	assert.Len(t, got, 2)
}

// ---- UpdateExpensesDetail ----

func TestRepo_UpdateExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`SELECT \* FROM "expenses_details" WHERE "expenses_details"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_id", "expenses_amount"}).AddRow(1, 10, 100))
	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 500))
	mk.ExpectExec(`UPDATE "expenses" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectExec(`UPDATE "expenses_details" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repository.NewExpensesRepository(db)
	got, err := r.UpdateExpensesDetail(testCtx, req.ExpensesDetail{
		ID:             1,
		ExpensesId:     10,
		ExpensesType:   "bill",
		ExpensesDesc:   "electricity",
		ExpensesAmount: 250,
	})

	assert.NoError(t, err)
	assert.Equal(t, "bill", got.ExpensesType)
	assert.Equal(t, float64(250), got.ExpensesAmount)
}

func TestRepo_UpdateExpensesDetail_DetailNotFound(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`SELECT \* FROM "expenses_details"`).
		WillReturnError(sql.ErrNoRows)
	mk.ExpectRollback()

	r := repository.NewExpensesRepository(db)
	got, err := r.UpdateExpensesDetail(testCtx, req.ExpensesDetail{ID: 999, ExpensesId: 10, ExpensesType: "x", ExpensesAmount: 1})

	assert.Error(t, err)
	assert.Equal(t, model.ExpensesDetail{}, got)
}

// ---- DeleteExpensesDetail ----

func TestRepo_DeleteExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`SELECT \* FROM "expenses_details" WHERE "expenses_details"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_id", "expenses_amount"}).AddRow(1, 10, 150))
	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 500))
	mk.ExpectExec(`UPDATE "expenses" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectExec(`DELETE FROM "expenses_details"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repository.NewExpensesRepository(db)
	ok, err := r.DeleteExpensesDetail(testCtx, req.DeleteExpensesDetailById{ID: 1})

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestRepo_DeleteExpensesDetail_NotFound(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`SELECT \* FROM "expenses_details"`).
		WillReturnError(sql.ErrNoRows)
	mk.ExpectRollback()

	r := repository.NewExpensesRepository(db)
	ok, err := r.DeleteExpensesDetail(testCtx, req.DeleteExpensesDetailById{ID: 999})

	assert.Error(t, err)
	assert.False(t, ok)
}
