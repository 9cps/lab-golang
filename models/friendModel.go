package models

import "gorm.io/gorm"

type Friend struct {
	gorm.Model
	F_NAME string
	L_NAME string
	TEL_NO string
}
