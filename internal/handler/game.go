/*package handler

import (
	"github.com/roma2099/uril-go/internal/database"
	"github.com/roma2099/uril-go/internal/model"
)
*/
// GetGamesHistory -> get 20 games acoding to params - info i wante to have is :gameplayer(you){gameID,postelo,seeds,starttime}[ordered by startTime];game{timePerPlayer, result}; gameplayer(adversery){userID};user(adversery){username}

// GetGame - >         game{timePerPlayer, result};plays{Sequence,Pit,Duration}[ordered by sequece,all with game id]  gameplayer(you){gameID,postelo,seeds,starttime}gameplayer(adversery){postelo,seeds,starttime};user(adversery){username}