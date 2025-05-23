package handler

import (
	"log"
	"strconv"

	"github.com/roma2099/uril-go/internal/database"
	"github.com/roma2099/uril-go/internal/model"

	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type GameHistoryItem struct {
	GameID        uint         `json:"game_id"`
	PostElo       uint16       `json:"post_elo"`
	Seeds         uint8        `json:"seeds"`
	StartTime     time.Time    `json:"start_time"`
	TimePerPlayer uint         `json:"time_per_player"`
	Result        model.Result `json:"result"`
	OpponentID    uint         `json:"opponent_id"`
	OpponentName  string       `json:"opponent_username"`
}
type GameOnline struct {
	Turn		 uint8        `json:"turn"`
	BoardState 		[12]uint          `json:"board_state"`
	Player1ID    uint         `json:"player1_id"`
	Player2ID    uint         `json:"player2_id"`
	Player1Conn  *websocket.Conn `json:"player1_conn"`
	Player2Conn  *websocket.Conn `json:"player2_conn"`
	Player1Elo   uint16       `json:"player1_elo"`
	Player2Elo   uint16       `json:"player2_elo"`
	Player1Seeds uint8        `json:"player1_seeds"`
	Player2Seeds uint8        `json:"player2_seeds"`
	Player1RemainTime time.Time `json:"player1_start_time"`
	Player2RemainTime time.Time `json:"player2_start_time"`
	
}
var GameOnlineList = make(map[uint]*GameOnline) // map[gameID]GameOnline

// GetGamesHistory -> get 20 games acoding to params - info i wante to have is :gameplayer(you){gameID,postelo,seeds,starttime}[ordered by startTime];game{timePerPlayer, result}; gameplayer(adversery){userID};user(adversery){username}
func GetGamesHistory(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	var history []GameHistoryItem
	db := database.DB
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
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
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "-", "errors": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Games found", "data": history})
}

func GetGame(c *fiber.Ctx) error {

	id := c.Params("id")
	var game model.Game
	db := database.DB
	err := db.Preload("Plays").
		Preload("GamePlayers.Player").
		First(&game, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Game not found", "errors": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Game found", "data": game})
}

// GetGame - >         game{timePerPlayer, result};plays{Sequence,Pit,Duration}[ordered by sequece,all with game userID gameplayer(you){gameID,postelo,seeds,starttime}gameplayer(adversery){postelo,seeds,starttime};user(adversery){username}

func GameWebSocketHandler(c *websocket.Conn) error {
	defer func() {
		c.Close()
	}()
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := uint(claims["user_id"].(float64))
	gameID := c.Params("id")
	var player_num uint8
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		log.Println("Invalid game ID:", err)
		return c.WriteMessage(websocket.TextMessage, []byte("Invalid game ID"))
	}
	gameIDUint := uint(gameIDInt)
	gameOnline, exists := GameOnlineList[gameIDUint]
	if !exists {
		log.Println("Game not found:", gameID)
		return c.WriteMessage(websocket.TextMessage, []byte("Game not found"))
	}
	if gameOnline.Player1ID != id && gameOnline.Player2ID != id {
		log.Println("User not in game:", id)
		return c.WriteMessage(websocket.TextMessage, []byte("User not in game"))
	}
	if gameOnline.Player1ID == id {
		gameOnline.Player1Conn = c
		player_num = 1
	} else {
		gameOnline.Player2Conn = c
		player_num = 2
	}
	// Send initial game state to the client
	err = c.WriteJSON(gameOnline)
	if err != nil {
		log.Println("Error sending initial game state:", err)
		return c.WriteMessage(websocket.TextMessage, []byte("Error sending initial game state"))
	}
	c.WriteMessage(websocket.TextMessage, []byte("Welcome to the game!"))
	
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error reading message:", err)
			}
			return nil
		}
		if msgType == websocket.TextMessage {
			if gameOnline.Turn != player_num {
				log.Println("It's not your turn")
				c.WriteMessage(websocket.TextMessage, []byte("It's not your turn"))
				continue
			}
			result,err:=gameLogic(gameOnline, msg, player_num)
			if err != nil {
				log.Println("Error in game logic:", err)
				c.WriteMessage(websocket.TextMessage, []byte("Error in game logic"))
				continue
			}

			gameOnline.Player1Conn.Conn.WriteJSON(gameOnline)
			gameOnline.Player2Conn.WriteJSON(gameOnline)
			if result == "GAME_OVER" {
				log.Println("Game Over")
				// Save game result to database
				//Remove game from GameOnlineList
				delete(GameOnlineList, gameIDUint)
			}
			
		} else {
			log.Println("websocket message received of type", msgType)
		}
	}

}
func gameLogic(gameOnline *GameOnline, msg []byte, player_num uint8) (string,error) {
	// Implement your game logic here
	// For example, you can update the game state based on the received message
	// and check if the game is over or if a player has won.
	
	// Example: Update the turn
	if player_num == 1 {
		gameOnline.Turn = 2
	} else {
		gameOnline.Turn = 1
	}

	return "CONTINUE",nil
}
