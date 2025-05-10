package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Winner string
const (
	Player1 Winner = "player1"
	Player2 Winner = "player2"
	Draw Winner = "draw"
)

func (w *Winner) Scan(value interface{}) error{
	str, ok:=value.(string)
	if !ok{
		return errors.New("invalid data for winner")
	}
	switch Winner(str){
	case Player1,Player2, Draw:
		*w = Winner(str)
		return nil
	default:
		return fmt.Errorf("invalid value for Winner: %s",str)
	}
}

func (w Winner) Value() (driver.Value, error){
	switch w{
	case Player1,Player2,Draw:
		return string(w),nil
	default :
		return nil, fmt.Errorf("invalid value for Winner: %s",w)
	}
}