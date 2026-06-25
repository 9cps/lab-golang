package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	req "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/request"
	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- mock service ----

type mockExpensesService struct {
	mock.Mock
}

func (m *mockExpensesService) InsertExpenses(_ context.Context, r req.Expenses) (res.ExpensesResponse, error) {
	args := m.Called(r)
	return args.Get(0).(res.ExpensesResponse), args.Error(1)
}
func (m *mockExpensesService) InsertExpensesDetail(_ context.Context, r req.ExpensesDetail) (res.ExpensesDetailResponse, error) {
	args := m.Called(r)
	return args.Get(0).(res.ExpensesDetailResponse), args.Error(1)
}
func (m *mockExpensesService) GetListMoneyCard(_ context.Context) (res.ExpensesCard, error) {
	args := m.Called()
	return args.Get(0).(res.ExpensesCard), args.Error(1)
}
func (m *mockExpensesService) GetListMoneyCardDetail(_ context.Context, r req.GetExpensesDetailById) ([]res.ExpensesDetailResponse, error) {
	args := m.Called(r)
	return args.Get(0).([]res.ExpensesDetailResponse), args.Error(1)
}
func (m *mockExpensesService) UpdateExpensesDetail(_ context.Context, r req.ExpensesDetail) (res.ExpensesDetailResponse, error) {
	args := m.Called(r)
	return args.Get(0).(res.ExpensesDetailResponse), args.Error(1)
}
func (m *mockExpensesService) DeleteExpensesDetail(_ context.Context, r req.DeleteExpensesDetailById) (bool, error) {
	args := m.Called(r)
	return args.Bool(0), args.Error(1)
}

// ---- helpers ----

func setupTest(m *mockExpensesService) (handler.ExpensesHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	h := handler.NewExpensesHandler(m)
	router := gin.New()
	return h, router
}

func performRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	r := httptest.NewRequest(method, path, &buf)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

// ---- CreateExpenses ----

func TestCreateExpenses_Success(t *testing.T) {
	m := new(mockExpensesService)
	body := req.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000}
	m.On("InsertExpenses", body).
		Return(res.ExpensesResponse{ExpensesMoney: 1000, ExpensesBalance: 1000}, nil)

	h, router := setupTest(m)
	router.PUT("/expenses", h.CreateExpenses)

	w := performRequest(router, http.MethodPut, "/expenses", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp res.DefaultResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, string(res.Success), resp.Status)
	m.AssertExpectations(t)
}

func TestCreateExpenses_InvalidJSON(t *testing.T) {
	m := new(mockExpensesService)
	h, router := setupTest(m)
	router.PUT("/expenses", h.CreateExpenses)

	r := httptest.NewRequest(http.MethodPut, "/expenses", bytes.NewBufferString("{bad json"))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "InsertExpenses")
}

func TestCreateExpenses_ServiceError(t *testing.T) {
	m := new(mockExpensesService)
	body := req.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000}
	m.On("InsertExpenses", body).
		Return(res.ExpensesResponse{}, errors.New("db error"))

	h, router := setupTest(m)
	router.PUT("/expenses", h.CreateExpenses)

	w := performRequest(router, http.MethodPut, "/expenses", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

// ---- UpsertExpensesDetail (create path: ID == 0) ----

func TestUpsertExpensesDetail_Create_Success(t *testing.T) {
	m := new(mockExpensesService)
	body := req.ExpensesDetail{ExpensesId: 1, ExpensesType: "food", ExpensesDesc: "lunch", ExpensesAmount: 120}
	m.On("InsertExpensesDetail", body).
		Return(res.ExpensesDetailResponse{
			ExpensesId:     body.ExpensesId,
			ExpensesType:   body.ExpensesType,
			ExpensesDesc:   body.ExpensesDesc,
			ExpensesAmount: body.ExpensesAmount,
		}, nil)

	h, router := setupTest(m)
	router.PUT("/expenses/details", h.UpsertExpensesDetail)

	w := performRequest(router, http.MethodPut, "/expenses/details", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "UpdateExpensesDetail")
}

func TestUpsertExpensesDetail_Create_ServiceError(t *testing.T) {
	m := new(mockExpensesService)
	body := req.ExpensesDetail{ExpensesId: 1, ExpensesType: "food", ExpensesAmount: 120}
	m.On("InsertExpensesDetail", body).
		Return(res.ExpensesDetailResponse{}, errors.New("db error"))

	h, router := setupTest(m)
	router.PUT("/expenses/details", h.UpsertExpensesDetail)

	w := performRequest(router, http.MethodPut, "/expenses/details", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

// ---- UpsertExpensesDetail (update path: ID != 0) ----

func TestUpsertExpensesDetail_Update_Success(t *testing.T) {
	m := new(mockExpensesService)
	body := req.ExpensesDetail{ID: 5, ExpensesId: 1, ExpensesType: "food", ExpensesDesc: "dinner", ExpensesAmount: 250}
	m.On("UpdateExpensesDetail", body).
		Return(res.ExpensesDetailResponse{
			ID:             body.ID,
			ExpensesId:     body.ExpensesId,
			ExpensesType:   body.ExpensesType,
			ExpensesDesc:   body.ExpensesDesc,
			ExpensesAmount: body.ExpensesAmount,
		}, nil)

	h, router := setupTest(m)
	router.PUT("/expenses/details", h.UpsertExpensesDetail)

	w := performRequest(router, http.MethodPut, "/expenses/details", body)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "InsertExpensesDetail")
}

func TestUpsertExpensesDetail_Update_ServiceError(t *testing.T) {
	m := new(mockExpensesService)
	body := req.ExpensesDetail{ID: 5, ExpensesId: 1, ExpensesType: "food", ExpensesAmount: 250}
	m.On("UpdateExpensesDetail", body).
		Return(res.ExpensesDetailResponse{}, errors.New("not found"))

	h, router := setupTest(m)
	router.PUT("/expenses/details", h.UpsertExpensesDetail)

	w := performRequest(router, http.MethodPut, "/expenses/details", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestUpsertExpensesDetail_InvalidJSON(t *testing.T) {
	m := new(mockExpensesService)
	h, router := setupTest(m)
	router.PUT("/expenses/details", h.UpsertExpensesDetail)

	r := httptest.NewRequest(http.MethodPut, "/expenses/details", bytes.NewBufferString("{"))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "InsertExpensesDetail")
	m.AssertNotCalled(t, "UpdateExpensesDetail")
}

// ---- GetListMoneyCard ----

func TestGetListMoneyCard_Success(t *testing.T) {
	m := new(mockExpensesService)
	m.On("GetListMoneyCard").
		Return(res.ExpensesCard{TotalBalance: 500}, nil)

	h, router := setupTest(m)
	router.GET("/expenses", h.GetListMoneyCard)

	w := performRequest(router, http.MethodGet, "/expenses", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertCalled(t, "GetListMoneyCard")
}

// ---- GetListMoneyCardDetail ----

func TestGetListMoneyCardDetail_Success(t *testing.T) {
	m := new(mockExpensesService)
	r := req.GetExpensesDetailById{Id: 10}
	m.On("GetListMoneyCardDetail", r).
		Return([]res.ExpensesDetailResponse{{ExpensesId: 10, ExpensesType: "food"}}, nil)

	h, router := setupTest(m)
	router.POST("/expenses/details", h.GetListMoneyCardDetail)

	w := performRequest(router, http.MethodPost, "/expenses/details", r)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestGetListMoneyCardDetail_MissingBody(t *testing.T) {
	m := new(mockExpensesService)
	h, router := setupTest(m)
	router.POST("/expenses/details", h.GetListMoneyCardDetail)

	w := performRequest(router, http.MethodPost, "/expenses/details", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "GetListMoneyCardDetail")
}

// ---- DeleteExpensesDetail ----

func TestDeleteExpensesDetail_Success(t *testing.T) {
	m := new(mockExpensesService)
	r := req.DeleteExpensesDetailById{ID: 1, ExpensesId: 10}
	m.On("DeleteExpensesDetail", r).Return(true, nil)

	h, router := setupTest(m)
	router.DELETE("/expenses/details", h.DeleteExpensesDetail)

	w := performRequest(router, http.MethodDelete, "/expenses/details", r)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteExpensesDetail_ServiceError(t *testing.T) {
	m := new(mockExpensesService)
	r := req.DeleteExpensesDetailById{ID: 1, ExpensesId: 10}
	m.On("DeleteExpensesDetail", r).Return(false, errors.New("not found"))

	h, router := setupTest(m)
	router.DELETE("/expenses/details", h.DeleteExpensesDetail)

	w := performRequest(router, http.MethodDelete, "/expenses/details", r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteExpensesDetail_MissingBody(t *testing.T) {
	m := new(mockExpensesService)
	h, router := setupTest(m)
	router.DELETE("/expenses/details", h.DeleteExpensesDetail)

	w := performRequest(router, http.MethodDelete, "/expenses/details", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "DeleteExpensesDetail")
}
