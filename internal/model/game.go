package model

import(
	"gorm.io/gorm"
	"time"
)

type Game struct{
	gorm.Model
	TimePerPlayer uint `gorm:"not null" validate:"required" json:"time_per_player"`
	Duration 	time.Duration `gorm:"not null" validate:"required" json:"duration"`
	Result		Result `gorm:"not null;type:text" validate:"required,oneof=timeout gameover resign" json:"result"`
	Plays 		[]Play `json:"plays"`
}
