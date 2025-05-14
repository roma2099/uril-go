package handler

import (
	"strconv"

	"github.com/roma2099/uril-go/internal/database"
	"github.com/roma2099/uril-go/internal/model"
	
	"github.com/gofiber/fiber/v2"
	"time"

)
type GameHistoryItem struct {
	GameID         uint          `json:"game_id"`
	PostElo        uint16        `json:"post_elo"`
	Seeds          uint8         `json:"seeds"`
	StartTime      time.Time     `json:"start_time"`
	TimePerPlayer  uint          `json:"time_per_player"`
	Result         model.Result  `json:"result"`
	OpponentID     uint          `json:"opponent_id"`
	OpponentName   string        `json:"opponent_username"`
}



// GetGamesHistory -> get 20 games acoding to params - info i wante to have is :gameplayer(you){gameID,postelo,seeds,starttime}[ordered by startTime];game{timePerPlayer, result}; gameplayer(adversery){userID};user(adversery){username}
func GetGamesHistory(c *fiber.Ctx) error{
	userID := c.Params("user_id")
	pageStr := c.Query("page","1")
	limitStr := c.Query("limit","20")
	var history []GameHistoryItem
	db:=database.DB
	page,err:=strconv.Atoi(pageStr)
	if err!=nil || page<1{
		page=1
	}
	limit,err :=strconv.Atoi(limitStr)
	if err!=nil || limit<1{
		limit=20
	}
	
	offset := (page - 1) * limit
	err = db.Table("game_players AS gp1").
		Select(`
			gp1.game_id, gp1.post_elo, gp1.seeds, gp1.start_time,
			g.time_per_player, g.result,
			gp2.user_id AS opponent_id,
			u.username AS opponent_name
		`).
		Joins("JOIN games g ON g.userID gp1.game_id").
		Joins("JOIN game_players gp2 ON gp2.game_id = gp1.game_id AND gp2.user_id != gp1.user_id").
		Joins("JOIN users u ON u.userID gp2.user_id").
		Where("gp1.user_id = ?", userID).
		Order("gp1.start_time DESC").
		Limit(limit).
		Offset(offset).
		Scan(&history).Error
	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "-", "errors": err.Error()})
	}
	return c.JSON(fiber.Map{"staus":"success","message":"Games found","data":history})
}

func GetGame(c *fiber.Ctx) error{

	id:=c.Params("id")
	var game model.Game
	db:=database.DB
	err	:=	db.Preload("Plays").
		Preload("GamePlayers.Player").
		First(&game,id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Game not found", "errors": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Game found", "data": game})
}
// GetGame - >         game{timePerPlayer, result};plays{Sequence,Pit,Duration}[ordered by sequece,all with game userID gameplayer(you){gameID,postelo,seeds,starttime}gameplayer(adversery){postelo,seeds,starttime};user(adversery){username}