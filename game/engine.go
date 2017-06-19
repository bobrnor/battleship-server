package game

import (
	"git.nulana.com/bobrnor/battleship-server/db"
)

type TurnResult uint8

type Engine struct{}

const (
	TurnResultMiss = iota
	TurnResultHit
)

var (
	engine *Engine
)

func DefaultEngine() *Engine {
	if engine == nil {
		engine = &Engine{}
	}
	return engine
}

func (e *Engine) SetGrid(room *db.Room, client *db.Client, gridData [13]uint8) error {
	return nil
}

func (e *Engine) Turn(room *db.Room, client *db.Client, x, y uint) (TurnResult, error) {
	return TurnResultMiss, nil
}
