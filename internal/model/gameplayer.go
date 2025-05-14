package model

import(
	"time"
)

type GamePlayer struct{
	UserID uint `gorm:"primaryKey;not null;index" validate:"required" json:"user_id"`
	Player User `gorm:"foreignKey:UserID"`
	GameID uint	`gorm:"primaryKey;not null;index" validate:"required" json:"game_id"`
	PlayerNum uint8 `gorm:"not null;check:player_num<=2" validate:"required,lte=2" json:"player_num"`
	PostElo uint16 `gorm:"not null;check:post_elo<=10000" validate:"required,lte=10000" json:"post_elo"`
	Seeds uint8 `gorm:"not null;check:seeds<=48" validate:"required,lte=48" json:"seeds"`
	StartTime time.Time `gorm:"not null;index" validate:"required" json:"start_time"`
}
// Depending on the DB may need to run this on the to create composite primarykey: CREATE UNIQUE INDEX IF NOT EXISTS idx_user_game ON game_players(user_id, game_id);
