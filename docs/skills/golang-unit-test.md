# Skill: Unit Test in Go — Best Practices (Table-Driven Pattern)

อ้างอิง: https://medium.com/@danarcahyaa/unit-test-in-go-best-practices-for-easier-testing-79d194fe9a54

---

## แนวคิดหลัก

แทนที่จะสร้าง test function แยกสำหรับแต่ละ scenario ให้ใช้ **Table-Driven Test** โดยสร้าง slice of struct ที่กำหนด input, expected output และ mock function สำหรับแต่ละ case ไว้ในที่เดียว

ข้อดี:
- ไม่ต้องสร้าง test function ใหม่ทุกครั้ง เพียงเพิ่ม entry ใน slice
- โค้ดกระชับ อ่านง่าย จัดการง่าย
- ลด boilerplate ซ้ำซ้อน

---

## โครงสร้างไฟล์

```
test/
  init.go                    # centralize mock + service initialization
  product_services_test.go   # test cases
  mocks/
    product_repository_mock.go
```

---

## 1. `test/init.go` — Centralize Configuration

สร้างไฟล์ `init.go` ไว้ใน package `test` เพื่อ initialize mock และ service ร่วมกัน

```go
// test/init.go
package test

import (
    "yourmodule/services"
    "yourmodule/test/mocks"
)

var productRepositoryMock *mocks.ProductRepositoryMock
var productServices *services.ProductService

func init() {
    productRepositoryMock = mocks.NewProductRepositoryMock()
    productServices = services.NewProductService(productRepositoryMock)
}
```

---

## 2. Mock Struct

```go
// test/mocks/product_repository_mock.go
package mocks

import (
    "yourmodule/entity"
    "github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
    Mock mock.Mock
}

func NewProductRepositoryMock() *ProductRepositoryMock {
    return &ProductRepositoryMock{}
}

func (p *ProductRepositoryMock) CreateOne(product *entity.Product) error {
    args := p.Mock.Called(product)
    return args.Error(0)
}
```

---

## 3. Table-Driven Test Pattern

```go
// test/product_services_test.go
package test

import (
    "errors"
    "yourmodule/entity"
    "yourmodule/models"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
    // struct กำหนด schema ของแต่ละ test case
    type CreateProductTest struct {
        TestName       string
        Request        *models.ProductRequest
        ExpectedError  error
        ExpectedResult *models.ProductResponse
        Mock           func()            // mock setup เฉพาะของแต่ละ case
    }

    tests := []CreateProductTest{
        // Test case 1: Success
        {
            TestName: "Shouldn't return an error",
            Request: &models.ProductRequest{
                Name:        "product",
                Price:       120,
                Description: "Description",
                Stock:       10,
                Category:    "category",
            },
            ExpectedError: nil,
            ExpectedResult: &models.ProductResponse{
                Name:        "product",
                Price:       120,
                Description: "Description",
                Stock:       10,
                Category:    "category",
            },
            Mock: func() {
                productMock := &entity.Product{
                    Name: "product", Price: 120,
                    Description: "Description", Stock: 10, Category: "category",
                }
                productRepositoryMock.Mock.On("CreateOne", productMock).Return(nil)
            },
        },
        // Test case 2: Error
        {
            TestName:       "Should return an error",
            Request:        &models.ProductRequest{Name: "product", Price: 120},
            ExpectedError:  errors.New("failed to create product"),
            ExpectedResult: nil,
            Mock: func() {
                productMock := &entity.Product{Name: "product", Price: 120}
                productRepositoryMock.Mock.On("CreateOne", productMock).
                    Return(errors.New("failed to create product"))
            },
        },
    }

    // วน loop รัน test แต่ละ case
    for _, test := range tests {
        t.Run(test.TestName, func(t *testing.T) {
            test.Mock()    // setup mock สำหรับ case นี้

            result, err := productServices.CreateProduct(test.Request)

            assert.Equal(t, test.ExpectedError, err)
            assert.Equal(t, test.ExpectedResult, result)
        })
    }
}
```

---

## 4. นำไปใช้กับโปรเจกต์นี้ (3-Layer)

### Controller Layer

```go
type CreateExpensesTest struct {
    TestName       string
    RequestBody    interface{}
    ExpectedStatus int
    ExpectedMsg    string
    Mock           func(m *mockExpensesServices)
}

tests := []CreateExpensesTest{
    {
        TestName:       "Success",
        RequestBody:    req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000},
        ExpectedStatus: http.StatusOK,
        ExpectedMsg:    string(res_dtos.Success),
        Mock: func(m *mockExpensesServices) {
            body := req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000}
            m.On("InsertExpenses", body).Return(models.Expenses{ExpensesMoney: 1000})
        },
    },
    {
        TestName:       "InvalidJSON",
        RequestBody:    nil,
        ExpectedStatus: http.StatusBadRequest,
        Mock:           func(m *mockExpensesServices) {},
    },
}

for _, tc := range tests {
    t.Run(tc.TestName, func(t *testing.T) {
        m := new(mockExpensesServices)
        tc.Mock(m)
        controller, router := setupTest(m)
        router.PUT("/Expenses/CreateExpenses", controller.CreateExpenses)
        w := performRequest(router, http.MethodPut, "/Expenses/CreateExpenses", tc.RequestBody)
        assert.Equal(t, tc.ExpectedStatus, w.Code)
        m.AssertExpectations(t)
    })
}
```

### Service Layer

```go
type InsertExpensesTest struct {
    TestName       string
    Request        req_dtos.Expenses
    ExpectedMoney  float32
    MockDB         func(mk sqlmock.Sqlmock)
    MockRepo       func(r *mockExpensesRepo)
}

tests := []InsertExpensesTest{
    {
        TestName:      "BelowThreshold_Success",
        Request:       req_dtos.Expenses{ExpensesMonth: 4, ExpensesYear: 2026, ExpensesMoney: 1000},
        ExpectedMoney: 1000,
        MockDB: func(mk sqlmock.Sqlmock) {
            mk.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "expenses"`)).
                WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
        },
        MockRepo: func(r *mockExpensesRepo) {
            r.On("InsertExpenses", mock.AnythingOfType("models.Expenses")).
                Return(models.Expenses{ExpensesMoney: 1000})
        },
    },
}

for _, tc := range tests {
    t.Run(tc.TestName, func(t *testing.T) {
        mk, cleanup := setupMockDB(t)
        defer cleanup()
        tc.MockDB(mk)
        repo := new(mockExpensesRepo)
        tc.MockRepo(repo)
        svc := services.NewExpensesServiceImpl(repo)
        result := svc.InsertExpenses(tc.Request)
        assert.Equal(t, tc.ExpectedMoney, result.ExpensesMoney)
        repo.AssertExpectations(t)
    })
}
```

---

## Best Practices สรุป

| แนวทาง | รายละเอียด |
|---|---|
| `init.go` | Centralize การ initialize mock และ service ใน test package |
| Table-Driven Test | ใช้ struct slice แทนการสร้าง test function แยกทีละ case |
| `Mock func()` field | แต่ละ test case กำหนด mock expectation ของตัวเอง |
| `t.Run(name, func)` | รันแต่ละ case เป็น subtest ที่มีชื่อชัดเจน |
| ชื่อ test case | ใช้ภาษาธรรมชาติอธิบาย scenario เช่น "Shouldn't return an error" |
| `AssertExpectations` | ยืนยันว่า mock ถูกเรียกครบตาม expectation ทุกครั้ง |

---

> Source: https://medium.com/@danarcahyaa/unit-test-in-go-best-practices-for-easier-testing-79d194fe9a54  
> Source code: https://github.com/DanarCahyadi12/go-unit-test
