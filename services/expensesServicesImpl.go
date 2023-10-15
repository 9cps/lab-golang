package services

import (
	req_dtos "github.com/9cps/api-go-gin/dtos/request"
	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/models"
	"github.com/9cps/api-go-gin/repository"
)

type ExpensesServiceImpl struct {
	ExpensesRopository repository.ExpensesRopository
}

func NewExpensesServiceImpl(expensesRopository repository.ExpensesRopository) ExpensesServices {
	return &ExpensesServiceImpl{
		ExpensesRopository: expensesRopository,
	}
}

func (s *ExpensesServiceImpl) InsertExpenses(req req_dtos.Expenses) models.Expenses {

	isCondition := DeleteExpensesIfCountExceeds(req, 6)

	if isCondition {
		obj := models.Expenses{
			ExpensesMonth:   req.ExpensesMonth,
			ExpensesYear:    req.ExpensesYear,
			ExpensesMoney:   req.ExpensesMoney, // จำนวนเงินตั้งต้น
			ExpensesBalance: req.ExpensesMoney, // จำนวนเงินคงเหลือ
		}

		result := s.ExpensesRopository.InsertExpenses(obj)
		return result
	}
	return models.Expenses{}
}

func (s *ExpensesServiceImpl) InsertExpensesDetail(req req_dtos.ExpensesDetail) models.ExpensesDetail {

	obj := models.ExpensesDetail{
		ExpensesId:     req.ExpensesId,
		ExpensesType:   req.ExpensesType,
		ExpensesDesc:   req.ExpensesDesc,
		ExpensesAmount: req.ExpensesAmount,
	}

	// Create the ExpensesDetail record
	result := s.ExpensesRopository.InsertExpensesDetail(obj)
	return result
}

func DeleteExpensesIfCountExceeds(req req_dtos.Expenses, threshold int64) bool {
	// Count the number of records
	var count int64
	if err := initializers.DB.Model(&models.Expenses{}).Count(&count).Error; err != nil {
		return false
	}

	if count >= threshold {
		// Find and delete the record with the smallest ID
		var expenses models.Expenses
		if err := initializers.DB.Order("id").First(&expenses).Error; err != nil {
			return false
		}

		// ใช้ลบ Unscoped เนื่องจากพื้น db มีจำกัด
		if err := initializers.DB.Unscoped().Delete(&expenses).Error; err != nil {
			return false
		}
	}
	return true
}

func (s *ExpensesServiceImpl) GetListMoneyCard() res_dtos.ExpensesCard {
	result := s.ExpensesRopository.GetListMoneyCard()
	return result
}

func (s *ExpensesServiceImpl) GetListMoneyCardDetail(req req_dtos.GetExpensesDetailById) []models.ExpensesDetail {
	result := s.ExpensesRopository.GetListMoneyCardDetail(req)
	return result
}
