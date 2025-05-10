package model

import(
	"gorm.io/gorm"
	"time"
)

type Game struct{
	gorm.Model
	Player1Ref	uint `gorm:"not null;index" validate:"required" json:"player1"`
	Player1		User `gorm:"foreignKey:Player1Ref"`
	Player2Ref	uint `gorm:"not null;index" validate:"required" json:"player2"`
	Player2		User `gorm:"foreignKey:Player2Ref"`
	Winner	 	Winner `gorm:"not null;type:text" validate:"required,oneof=player1 player2 draw" json:"winner"`	
	TimePerPlayer uint `gorm:"not null" validate:"required" json:"time_per_player"`
	StartTime 	time.Time `gorm:"not null" validate:"required" json:"start_time"`
	EndTime 	time.Time `gorm:"not null" validate:"required" json:"end_time"`
	P1Seeds	uint `gorm:"not null;check:p1_seeds<=48" validate:"required,lte=48" json:"p1_seeds"`
	P2Seeds	uint `gorm:"not null;check:p2_seeds<=48" validate:"required,lte=48" json:"p2_seeds"`
	P1Elo 		uint `gorm:"not null;check:p1_elo<=10000" validate:"required,lte=10000" json:"p1_elo"`
	P2Elo 		uint `gorm:"not null;check:p2_elo<=10000" validate:"required,lte=10000" json:"p2_elo"`
	Result		Result `gorm:"not null;type:text" validate:"required,oneof=time\\ out game\\ over resign" json:"result"`
	Plays 		[]Play
}

