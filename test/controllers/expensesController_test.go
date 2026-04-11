package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/9cps/api-go-gin/controllers"
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExpensesServices struct {
	mock.Mock
}

func (m *mockExpensesServices) InsertExpenses(req req_dtos.Expenses) models.Expenses {
	args := m.Called(req)
	return args.Get(0).(models.Expenses)
}
func (m *mockExpensesServices) InsertExpensesDetail(req req_dtos.ExpensesDetail) models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).(models.ExpensesDetail)
}
func (m *mockExpensesServices) GetListMoneyCard() res_dtos.ExpensesCard {
	args := m.Called()
	return args.Get(0).(res_dtos.ExpensesCard)
}
func (m *mockExpensesServices) GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).([]models.ExpensesDetail)
}
func (m *mockExpensesServices) UpdateExpensesDetail(req req_dtos.UpdateExpensesDetail) models.ExpensesDetail {
	args := m.Called(req)
	return args.Get(0).(models.ExpensesDetail)
}
func (m *mockExpensesServices) DeleteExpensesDetail(req req_dtos.DeleteExpensesDetailById) bool {
	args := m.Called(req)
	return args.Bool(0)
}

func setupTest(m *mockExpensesServices) (*controllers.ExpensesController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	controller := controllers.NewExpensesController(m)
	router := gin.New()
	return controller, router
}

func performRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ---------------- CreateExpenses ----------------

func TestCreateExpenses_Success(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000}
	m.On("InsertExpenses", body).Return(models.Expenses{
		ExpensesMonth:   body.ExpensesMonth,
		ExpensesYear:    body.ExpensesYear,
		ExpensesMoney:   body.ExpensesMoney,
		ExpensesBalance: body.ExpensesMoney,
	})

	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpenses", controller.CreateExpenses)

	w := performRequest(router, http.MethodPut, "/Expenses/CreateExpenses", body)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp res_dtos.DefaultResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, string(res_dtos.Success), resp.Status)
	m.AssertExpectations(t)
}

func TestCreateExpenses_InvalidJSON(t *testing.T) {
	m := new(mockExpensesServices)
	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpenses", controller.CreateExpenses)

	req := httptest.NewRequest(http.MethodPut, "/Expenses/CreateExpenses", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "InsertExpenses")
}

func TestCreateExpenses_ServiceReturnsEmpty(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000}
	m.On("InsertExpenses", body).Return(models.Expenses{})

	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpenses", controller.CreateExpenses)

	w := performRequest(router, http.MethodPut, "/Expenses/CreateExpenses", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

// ---------------- CreateExpensesDetail ----------------

func TestCreateExpensesDetail_Success(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.ExpensesDetail{ExpensesId: 1, ExpensesType: "food", ExpensesDesc: "lunch", ExpensesAmount: 120}
	m.On("InsertExpensesDetail", body).Return(models.ExpensesDetail{
		ExpensesId:     body.ExpensesId,
		ExpensesType:   body.ExpensesType,
		ExpensesDesc:   body.ExpensesDesc,
		ExpensesAmount: body.ExpensesAmount,
	})

	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpensesDetail", controller.CreateExpensesDetail)

	w := performRequest(router, http.MethodPut, "/Expenses/CreateExpensesDetail", body)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestCreateExpensesDetail_ServiceReturnsEmpty(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.ExpensesDetail{ExpensesId: 1}
	m.On("InsertExpensesDetail", body).Return(models.ExpensesDetail{})

	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpensesDetail", controller.CreateExpensesDetail)

	w := performRequest(router, http.MethodPut, "/Expenses/CreateExpensesDetail", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

func TestCreateExpensesDetail_InvalidJSON(t *testing.T) {
	m := new(mockExpensesServices)
	controller, router := setupTest(m)
	router.PUT("/Expenses/CreateExpensesDetail", controller.CreateExpensesDetail)

	req := httptest.NewRequest(http.MethodPut, "/Expenses/CreateExpensesDetail", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "InsertExpensesDetail")
}

// ---------------- GetListMoneyCard ----------------

func TestGetListMoneyCard_Success(t *testing.T) {
	m := new(mockExpensesServices)
	m.On("GetListMoneyCard").Return(res_dtos.ExpensesCard{TotalBalance: 500, PercentBalance: 50})

	controller, router := setupTest(m)
	router.GET("/Expenses/GetListMoneyCard", controller.GetListMoneyCard)

	w := performRequest(router, http.MethodGet, "/Expenses/GetListMoneyCard", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertCalled(t, "GetListMoneyCard")
}

// ---------------- GetListMoneyCardDetail ----------------

func TestGetListMoneyCardDetail_Success(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.GetExpensesDetailById{Id: 10}
	m.On("GetListMoneyCardDetail", body).
		Return([]models.ExpensesDetail{{ExpensesId: 10, ExpensesType: "food"}})

	controller, router := setupTest(m)
	router.POST("/Expenses/GetListMoneyCardDetail", controller.GetListMoneyCardDetail)

	w := performRequest(router, http.MethodPost, "/Expenses/GetListMoneyCardDetail", body)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestGetListMoneyCardDetail_InvalidJSON(t *testing.T) {
	m := new(mockExpensesServices)
	controller, router := setupTest(m)
	router.POST("/Expenses/GetListMoneyCardDetail", controller.GetListMoneyCardDetail)

	req := httptest.NewRequest(http.MethodPost, "/Expenses/GetListMoneyCardDetail", bytes.NewBufferString("nope"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "GetListMoneyCardDetail")
}

// ---------------- UpdateExpensesDetail ----------------

func TestUpdateExpensesDetail_Success(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.UpdateExpensesDetail{Id: 1, ExpensesType: "food", ExpensesDesc: "dinner", ExpensesAmount: 250}
	m.On("UpdateExpensesDetail", body).Return(models.ExpensesDetail{
		ExpensesType:   body.ExpensesType,
		ExpensesDesc:   body.ExpensesDesc,
		ExpensesAmount: body.ExpensesAmount,
	})

	controller, router := setupTest(m)
	router.PUT("/Expenses/UpdateExpensesDetail", controller.UpdateExpensesDetail)

	w := performRequest(router, http.MethodPut, "/Expenses/UpdateExpensesDetail", body)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestUpdateExpensesDetail_ServiceReturnsEmpty(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.UpdateExpensesDetail{Id: 1}
	m.On("UpdateExpensesDetail", body).Return(models.ExpensesDetail{})

	controller, router := setupTest(m)
	router.PUT("/Expenses/UpdateExpensesDetail", controller.UpdateExpensesDetail)

	w := performRequest(router, http.MethodPut, "/Expenses/UpdateExpensesDetail", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

func TestUpdateExpensesDetail_InvalidJSON(t *testing.T) {
	m := new(mockExpensesServices)
	controller, router := setupTest(m)
	router.PUT("/Expenses/UpdateExpensesDetail", controller.UpdateExpensesDetail)

	req := httptest.NewRequest(http.MethodPut, "/Expenses/UpdateExpensesDetail", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "UpdateExpensesDetail")
}

// ---------------- DeleteExpensesDetail ----------------

func TestDeleteExpensesDetail_Success(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.DeleteExpensesDetailById{Id: 1}
	m.On("DeleteExpensesDetail", body).Return(true)

	controller, router := setupTest(m)
	router.DELETE("/Expenses/DeleteExpensesDetail", controller.DeleteExpensesDetail)

	w := performRequest(router, http.MethodDelete, "/Expenses/DeleteExpensesDetail", body)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteExpensesDetail_ServiceReturnsFalse(t *testing.T) {
	m := new(mockExpensesServices)
	body := req_dtos.DeleteExpensesDetailById{Id: 1}
	m.On("DeleteExpensesDetail", body).Return(false)

	controller, router := setupTest(m)
	router.DELETE("/Expenses/DeleteExpensesDetail", controller.DeleteExpensesDetail)

	w := performRequest(router, http.MethodDelete, "/Expenses/DeleteExpensesDetail", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteExpensesDetail_InvalidJSON(t *testing.T) {
	m := new(mockExpensesServices)
	controller, router := setupTest(m)
	router.DELETE("/Expenses/DeleteExpensesDetail", controller.DeleteExpensesDetail)

	req := httptest.NewRequest(http.MethodDelete, "/Expenses/DeleteExpensesDetail", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "DeleteExpensesDetail")
}
