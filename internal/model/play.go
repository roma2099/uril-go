package model

import (
	"time"
)

type Play struct {
	GameID   uint      `gorm:"primaryKey;not null;index" validate:"required" json:"game_id"`
	Sequence uint8     `gorm:"primaryKey;not null" validate:"required" json:"sequence"`
	Pit      uint8     `gorm:"not null;check:pit<=11" validate:"required,lte=11" json:"pit"`
	Duration     time.Duration `gorm:"not null" validate:"required" json:"duration"`
	Game     Game      `gorm:"foreignKey:GameID;references:ID;constraint:OnDelete:CASCADE;" validate:"required" json:"-"`
}
