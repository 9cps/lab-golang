package service_test

import (
	"context"
	"errors"
	"testing"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/model"
	svc "github.com/9cps/api-go-gin/services/expenses-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- mock repository ----

type mockExpensesRepo struct {
	mock.Mock
}

func (m *mockExpensesRepo) InsertExpenses(ctx context.Context, obj model.Expenses) (model.Expenses, error) {
	args := m.Called(ctx, obj)
	return args.Get(0).(model.Expenses), args.Error(1)
}
func (m *mockExpensesRepo) InsertExpensesDetail(ctx context.Context, obj model.ExpensesDetail) (model.ExpensesDetail, error) {
	args := m.Called(ctx, obj)
	return args.Get(0).(model.ExpensesDetail), args.Error(1)
}
func (m *mockExpensesRepo) GetListMoneyCard(ctx context.Context) ([]model.Expenses, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Expenses), args.Error(1)
}
func (m *mockExpensesRepo) GetListMoneyCardDetail(ctx context.Context, r req.GetExpensesDetailById) ([]model.ExpensesDetail, error) {
	args := m.Called(ctx, r)
	return args.Get(0).([]model.ExpensesDetail), args.Error(1)
}
func (m *mockExpensesRepo) UpdateExpensesDetail(ctx context.Context, r req.ExpensesDetail) (model.ExpensesDetail, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(model.ExpensesDetail), args.Error(1)
}
func (m *mockExpensesRepo) DeleteExpensesDetail(ctx context.Context, r req.DeleteExpensesDetailById) (bool, error) {
	args := m.Called(ctx, r)
	return args.Bool(0), args.Error(1)
}
func (m *mockExpensesRepo) CountExpenses(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockExpensesRepo) DeleteOldestExpenses(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

var testCtx = context.Background()

// ---- InsertExpenses ----

func TestService_InsertExpenses_BelowThreshold_Success(t *testing.T) {
	repo := new(mockExpensesRepo)
	repo.On("CountExpenses", testCtx).Return(int64(2), nil)
	repo.On("InsertExpenses", testCtx, mock.AnythingOfType("model.Expenses")).
		Return(model.Expenses{ExpensesMoney: 1000, ExpensesBalance: 1000}, nil)

	s := svc.NewExpensesService(repo)
	result, err := s.InsertExpenses(testCtx, req.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000})

	assert.NoError(t, err)
	assert.Equal(t, float64(1000), result.ExpensesMoney)
	assert.Equal(t, float64(1000), result.ExpensesBalance)
	repo.AssertExpectations(t)
}

func TestService_InsertExpenses_AtThreshold_DeletesOldest(t *testing.T) {
	repo := new(mockExpensesRepo)
	repo.On("CountExpenses", testCtx).Return(int64(6), nil)
	repo.On("DeleteOldestExpenses", testCtx).Return(nil)
	repo.On("InsertExpenses", testCtx, mock.AnythingOfType("model.Expenses")).
		Return(model.Expenses{ExpensesMoney: 1000, ExpensesBalance: 1000}, nil)

	s := svc.NewExpensesService(repo)
	result, err := s.InsertExpenses(testCtx, req.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000})

	assert.NoError(t, err)
	assert.Equal(t, float64(1000), result.ExpensesMoney)
	repo.AssertExpectations(t)
}

func TestService_InsertExpenses_CountQueryError(t *testing.T) {
	repo := new(mockExpensesRepo)
	repo.On("CountExpenses", testCtx).Return(int64(0), errors.New("db error"))

	s := svc.NewExpensesService(repo)
	result, err := s.InsertExpenses(testCtx, req.Expenses{ExpensesMoney: 1000})

	assert.Error(t, err)
	assert.Equal(t, res.ExpensesResponse{}, result)
	repo.AssertNotCalled(t, "InsertExpenses")
}

// ---- InsertExpensesDetail ----

func TestService_InsertExpensesDetail(t *testing.T) {
	obj := model.ExpensesDetail{
		ExpensesId:     1,
		ExpensesType:   "food",
		ExpensesDesc:   "lunch",
		ExpensesAmount: 120,
	}
	repo := new(mockExpensesRepo)
	repo.On("InsertExpensesDetail", testCtx, obj).Return(obj, nil)

	s := svc.NewExpensesService(repo)
	result, err := s.InsertExpensesDetail(testCtx, req.ExpensesDetail{
		ExpensesId:     1,
		ExpensesType:   "food",
		ExpensesDesc:   "lunch",
		ExpensesAmount: 120,
	})

	assert.NoError(t, err)
	assert.Equal(t, "food", result.ExpensesType)
	assert.Equal(t, float64(120), result.ExpensesAmount)
	repo.AssertExpectations(t)
}

// ---- GetListMoneyCard ----

func TestService_GetListMoneyCard(t *testing.T) {
	rows := []model.Expenses{
		{ExpensesMoney: 1000, ExpensesBalance: 800},
		{ExpensesMoney: 2000, ExpensesBalance: 1500},
	}
	repo := new(mockExpensesRepo)
	repo.On("GetListMoneyCard", testCtx).Return(rows, nil)

	s := svc.NewExpensesService(repo)
	got, err := s.GetListMoneyCard(testCtx)

	assert.NoError(t, err)
	assert.Equal(t, float64(2300), got.TotalBalance)
	assert.Len(t, got.Data, 2)
	assert.Equal(t, float64(200), got.Data[0].TotalSpending)
	repo.AssertExpectations(t)
}

// ---- GetListMoneyCardDetail ----

func TestService_GetListMoneyCardDetail(t *testing.T) {
	r := req.GetExpensesDetailById{Id: 7}
	repo := new(mockExpensesRepo)
	repo.On("GetListMoneyCardDetail", testCtx, r).
		Return([]model.ExpensesDetail{{ExpensesId: 7}}, nil)

	s := svc.NewExpensesService(repo)
	got, err := s.GetListMoneyCardDetail(testCtx, r)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int32(7), got[0].ExpensesId)
	repo.AssertExpectations(t)
}

// ---- UpdateExpensesDetail ----

func TestService_UpdateExpensesDetail(t *testing.T) {
	r := req.ExpensesDetail{ID: 1, ExpensesId: 10, ExpensesType: "bill", ExpensesAmount: 300}
	repo := new(mockExpensesRepo)
	repo.On("UpdateExpensesDetail", testCtx, r).Return(model.ExpensesDetail{
		ExpensesType:   r.ExpensesType,
		ExpensesAmount: r.ExpensesAmount,
	}, nil)

	s := svc.NewExpensesService(repo)
	got, err := s.UpdateExpensesDetail(testCtx, r)

	assert.NoError(t, err)
	assert.Equal(t, "bill", got.ExpensesType)
	assert.Equal(t, float64(300), got.ExpensesAmount)
	repo.AssertExpectations(t)
}

// ---- DeleteExpensesDetail ----

func TestService_DeleteExpensesDetail_True(t *testing.T) {
	r := req.DeleteExpensesDetailById{ID: 1}
	repo := new(mockExpensesRepo)
	repo.On("DeleteExpensesDetail", testCtx, r).Return(true, nil)

	s := svc.NewExpensesService(repo)
	ok, err := s.DeleteExpensesDetail(testCtx, r)

	assert.NoError(t, err)
	assert.True(t, ok)
	repo.AssertExpectations(t)
}

func TestService_DeleteExpensesDetail_Error(t *testing.T) {
	r := req.DeleteExpensesDetailById{ID: 1}
	repo := new(mockExpensesRepo)
	repo.On("DeleteExpensesDetail", testCtx, r).Return(false, errors.New("not found"))

	s := svc.NewExpensesService(repo)
	ok, err := s.DeleteExpensesDetail(testCtx, r)

	assert.Error(t, err)
	assert.False(t, ok)
	repo.AssertExpectations(t)
}
