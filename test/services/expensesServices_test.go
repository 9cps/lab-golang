package services_test

import (
	"database/sql"
	"regexp"
	"testing"

	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	services "github.com/9cps/api-go-gin/services/service"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockExpensesRepo struct {
	mock.Mock
}

func (m *mockExpensesRepo) InsertExpenses(req models.Expenses) models.Expenses {
	args := m.Called(req)
	return args.Get(0).(models.Expenses)
}
func (m *mockExpensesRepo) InsertExpensesDetail(req models.ExpensesDetail) models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).(models.ExpensesDetail)
}
func (m *mockExpensesRepo) GetListMoneyCard() res_dtos.ExpensesCard {
	args := m.Called()
	return args.Get(0).(res_dtos.ExpensesCard)
}
func (m *mockExpensesRepo) GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).([]models.ExpensesDetail)
}
func (m *mockExpensesRepo) UpdateExpensesDetail(req req_dtos.UpdateExpensesDetail) models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).(models.ExpensesDetail)
}
func (m *mockExpensesRepo) DeleteExpensesDetail(req req_dtos.DeleteExpensesDetailById) bool {
	args := m.Called(req)
	return args.Bool(0)
}

// setupMockDB wires initializers.DB to a sqlmock-backed gorm.DB so the
// package-level DB calls inside ExpensesServiceImpl can be controlled.
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()
	sqlDB, mk, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	require.NoError(t, err)
	prev := initializers.DB
	initializers.DB = gdb
	return mk, func() {
		initializers.DB = prev
		_ = sqlDB.Close()
	}
}

// ---------------- InsertExpenses ----------------

func TestService_InsertExpenses_BelowThreshold_Success(t *testing.T) {
	mk, cleanup := setupMockDB(t)
	defer cleanup()

	mk.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "expenses"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	repo := new(mockExpensesRepo)
	repo.On("InsertExpenses", mock.AnythingOfType("models.Expenses")).
		Return(models.Expenses{ExpensesMoney: 1000, ExpensesBalance: 1000})

	svc := services.NewExpensesServiceImpl(repo)
	result := svc.InsertExpenses(req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000})

	assert.Equal(t, float32(1000), result.ExpensesMoney)
	assert.Equal(t, float32(1000), result.ExpensesBalance)
	repo.AssertExpectations(t)
}

func TestService_InsertExpenses_CountQueryError(t *testing.T) {
	mk, cleanup := setupMockDB(t)
	defer cleanup()

	mk.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "expenses"`)).
		WillReturnError(sql.ErrConnDone)

	repo := new(mockExpensesRepo)
	svc := services.NewExpensesServiceImpl(repo)

	result := svc.InsertExpenses(req_dtos.Expenses{ExpensesMoney: 1000})

	assert.Equal(t, models.Expenses{}, result)
	repo.AssertNotCalled(t, "InsertExpenses")
}

// ---------------- InsertExpensesDetail ----------------

func TestService_InsertExpensesDetail(t *testing.T) {
	expected := models.ExpensesDetail{
		ExpensesId:     1,
		ExpensesType:   "food",
		ExpensesDesc:   "lunch",
		ExpensesAmount: 120,
	}
	repo := new(mockExpensesRepo)
	repo.On("InsertExpensesDetail", expected).Return(expected)

	svc := services.NewExpensesServiceImpl(repo)
	result := svc.InsertExpensesDetail(req_dtos.ExpensesDetail{
		ExpensesId:     1,
		ExpensesType:   "food",
		ExpensesDesc:   "lunch",
		ExpensesAmount: 120,
	})

	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// ---------------- GetListMoneyCard ----------------

func TestService_GetListMoneyCard(t *testing.T) {
	want := res_dtos.ExpensesCard{TotalBalance: 500}
	repo := new(mockExpensesRepo)
	repo.On("GetListMoneyCard").Return(want)

	svc := services.NewExpensesServiceImpl(repo)
	got := svc.GetListMoneyCard()

	assert.Equal(t, float32(500), got.TotalBalance)
	repo.AssertExpectations(t)
}

// ---------------- GetListMoneyCardDetail ----------------

func TestService_GetListMoneyCardDetail(t *testing.T) {
	req := req_dtos.GetExpensesDetailById{Id: 7}
	repo := new(mockExpensesRepo)
	repo.On("GetListMoneyCardDetail", req).
		Return([]models.ExpensesDetail{{ExpensesId: 7}})

	svc := services.NewExpensesServiceImpl(repo)
	got := svc.GetListMoneyCardDetail(req)

	assert.Len(t, got, 1)
	assert.Equal(t, int32(7), got[0].ExpensesId)
	repo.AssertExpectations(t)
}

// ---------------- UpdateExpensesDetail ----------------

func TestService_UpdateExpensesDetail(t *testing.T) {
	req := req_dtos.UpdateExpensesDetail{Id: 1, ExpensesType: "bill", ExpensesAmount: 300}
	repo := new(mockExpensesRepo)
	repo.On("UpdateExpensesDetail", req).Return(models.ExpensesDetail{
		ExpensesType:   req.ExpensesType,
		ExpensesAmount: req.ExpensesAmount,
	})

	svc := services.NewExpensesServiceImpl(repo)
	got := svc.UpdateExpensesDetail(req)

	assert.Equal(t, "bill", got.ExpensesType)
	assert.Equal(t, float32(300), got.ExpensesAmount)
	repo.AssertExpectations(t)
}

// ---------------- DeleteExpensesDetail ----------------

func TestService_DeleteExpensesDetail_True(t *testing.T) {
	req := req_dtos.DeleteExpensesDetailById{Id: 1}
	repo := new(mockExpensesRepo)
	repo.On("DeleteExpensesDetail", req).Return(true)

	svc := services.NewExpensesServiceImpl(repo)

	assert.True(t, svc.DeleteExpensesDetail(req))
	repo.AssertExpectations(t)
}

func TestService_DeleteExpensesDetail_False(t *testing.T) {
	req := req_dtos.DeleteExpensesDetailById{Id: 1}
	repo := new(mockExpensesRepo)
	repo.On("DeleteExpensesDetail", req).Return(false)

	svc := services.NewExpensesServiceImpl(repo)

	assert.False(t, svc.DeleteExpensesDetail(req))
	repo.AssertExpectations(t)
}

// ---------------- Intentionally failing case ----------------

// func TestService_InsertExpensesDetail_FailingCase(t *testing.T) {
// 	repo := new(mockExpensesRepo)
// 	repo.On("InsertExpensesDetail", mock.AnythingOfType("models.ExpensesDetail")).
// 		Return(models.ExpensesDetail{
// 			ExpensesId:     1,
// 			ExpensesType:   "food",
// 			ExpensesAmount: 100,
// 		})

// 	svc := services.NewExpensesServiceImpl(repo)
// 	got := svc.InsertExpensesDetail(req_dtos.ExpensesDetail{
// 		ExpensesId:     1,
// 		ExpensesType:   "food",
// 		ExpensesAmount: 100,
// 	})

// 	assert.Equal(t, float32(999), got.ExpensesAmount, "intentional failure")
// }
