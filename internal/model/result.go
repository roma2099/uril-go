package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Result string
const (
	TimeOut Result = "time out"
	GameOver Result = "game over"
	Resign Result = "resign"
)

func (r *Result) Scan(value interface{}) error{
	str, ok:=value.(string)
	if !ok{
		return errors.New("invalid data for winner")
	}
	switch Result(str){
	case TimeOut,GameOver, Resign:
		*r = Result(str)
		return nil
	default:
		return fmt.Errorf("invalid value for Result: %s",str)
	}
}

func (r Result) Value() (driver.Value, error){
	switch r{
	case TimeOut,GameOver,Resign:
		return string(r),nil
	default :
		return nil, fmt.Errorf("invalid value for Result: %s",r)
	}
}