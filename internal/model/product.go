package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title 		string `gorm:"not null" json:"title"`
	Descrition 	string `gorm:"not null" json:"descrition"`
	Amount 		int `gorm:"not null" json:"amount"`
}