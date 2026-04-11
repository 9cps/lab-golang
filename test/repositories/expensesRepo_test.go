package repositories_test

import (
	"database/sql"
	"regexp"
	"testing"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	"github.com/9cps/api-go-gin/models"
	repo "github.com/9cps/api-go-gin/repositories/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

// ---------------- InsertExpenses ----------------

func TestRepo_InsertExpenses_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mk.ExpectCommit()

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.InsertExpenses(models.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000, ExpensesBalance: 1000})

	assert.Equal(t, float32(1000), got.ExpensesMoney)
	assert.NoError(t, mk.ExpectationsWereMet())
}

func TestRepo_InsertExpenses_Error(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses"`).WillReturnError(sql.ErrConnDone)
	mk.ExpectRollback()

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.InsertExpenses(models.Expenses{ExpensesMoney: 1000})

	assert.Equal(t, models.Expenses{}, got)
}

// ---------------- InsertExpensesDetail ----------------

func TestRepo_InsertExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses_details"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mk.ExpectCommit()

	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 1000))

	mk.ExpectBegin()
	mk.ExpectExec(`UPDATE "expenses" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.InsertExpensesDetail(models.ExpensesDetail{ExpensesId: 10, ExpensesAmount: 200, ExpensesType: "food"})

	assert.Equal(t, int32(10), got.ExpensesId)
}

func TestRepo_InsertExpensesDetail_CreateError(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectBegin()
	mk.ExpectQuery(`INSERT INTO "expenses_details"`).WillReturnError(sql.ErrConnDone)
	mk.ExpectRollback()

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.InsertExpensesDetail(models.ExpensesDetail{ExpensesId: 10})

	assert.Equal(t, models.ExpensesDetail{}, got)
}

// ---------------- GetListMoneyCard ----------------

func TestRepo_GetListMoneyCard_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "expenses_month", "expenses_year", "expenses_money", "expenses_balance"}).
		AddRow(1, 1, 2026, 1000, 800).
		AddRow(2, 2, 2026, 2000, 1500)
	mk.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM expenses ORDER BY expenses_month ASC`)).WillReturnRows(rows)

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.GetListMoneyCard()

	assert.Equal(t, float32(2300), got.TotalBalance)
	assert.Len(t, got.Data, 2)
}

// ---------------- GetListMoneyCardDetail ----------------

func TestRepo_GetListMoneyCardDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "expenses_id", "expenses_type", "expenses_desc", "expenses_amount"}).
		AddRow(1, 5, "food", "lunch", 120).
		AddRow(2, 5, "bill", "water", 300)

	mk.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM expenses_details WHERE expenses_id = $1 ORDER BY created_at DESC`)).
		WithArgs(5).
		WillReturnRows(rows)

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.GetListMoneyCardDetail(req_dtos.GetExpensesDetailById{Id: 5})

	assert.Len(t, got, 2)
}

// ---------------- UpdateExpensesDetail ----------------

func TestRepo_UpdateExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectQuery(`SELECT \* FROM "expenses_details" WHERE "expenses_details"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_id", "expenses_amount"}).AddRow(1, 10, 100))

	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 500))

	mk.ExpectBegin()
	mk.ExpectExec(`UPDATE "expenses" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	mk.ExpectBegin()
	mk.ExpectExec(`UPDATE "expenses_details" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.UpdateExpensesDetail(req_dtos.UpdateExpensesDetail{
		Id:             1,
		ExpensesType:   "bill",
		ExpensesDesc:   "electricity",
		ExpensesAmount: 250,
	})

	assert.Equal(t, "bill", got.ExpensesType)
	assert.Equal(t, float32(250), got.ExpensesAmount)
}

func TestRepo_UpdateExpensesDetail_DetailNotFound(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectQuery(`SELECT \* FROM "expenses_details"`).
		WillReturnError(sql.ErrNoRows)

	r := repo.NewExpensesRepositoryImpl(db)
	got := r.UpdateExpensesDetail(req_dtos.UpdateExpensesDetail{Id: 999})

	assert.Equal(t, models.ExpensesDetail{}, got)
}

// ---------------- DeleteExpensesDetail ----------------

func TestRepo_DeleteExpensesDetail_Success(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectQuery(`SELECT \* FROM "expenses_details" WHERE "expenses_details"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_id", "expenses_amount"}).AddRow(1, 10, 150))

	mk.ExpectQuery(`SELECT \* FROM "expenses" WHERE "expenses"."id" = \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "expenses_balance"}).AddRow(10, 500))

	mk.ExpectBegin()
	mk.ExpectExec(`UPDATE "expenses" SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	mk.ExpectBegin()
	mk.ExpectExec(`DELETE FROM "expenses_details"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mk.ExpectCommit()

	r := repo.NewExpensesRepositoryImpl(db)
	assert.True(t, r.DeleteExpensesDetail(req_dtos.DeleteExpensesDetailById{Id: 1}))
}

func TestRepo_DeleteExpensesDetail_NotFound(t *testing.T) {
	db, mk, cleanup := newRepo(t)
	defer cleanup()

	mk.ExpectQuery(`SELECT \* FROM "expenses_details"`).
		WillReturnError(sql.ErrNoRows)

	r := repo.NewExpensesRepositoryImpl(db)
	assert.False(t, r.DeleteExpensesDetail(req_dtos.DeleteExpensesDetailById{Id: 999}))
}
