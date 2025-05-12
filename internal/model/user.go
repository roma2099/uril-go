package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model	// includes the ID
	Username	string 	`gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email 		string 	`gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password 	string 	`gorm:"not null" validate:"required,min=6,max=50" json:"password"`
	Elo			uint16	`gorm:"not null;default:1200;check:Elo<=10000" validate:"lte=10000" json:"elo"`//not sure if it is suposed to be like this "check:Elo<=10000" or like this "check:elo<=10000"
	Country		string	`gorm:"not null;size:2" validate:"len=2,required" json:"country"`
}