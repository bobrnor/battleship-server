package game

import (
	"git.nulana.com/bobrnor/battleship-server/db"
	"github.com/pkg/errors"
)

type TurnResult uint8

type Engine struct{}

const (
	TurnResultMiss = iota
	TurnResultHit
)

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) SetGrid(dbRoom *db.Room, client *db.Client, gridData [13]uint8) error {
	rooms := MainRooms()
	room := rooms.Room(dbRoom.UID)
	if room == nil {
		return errors.Errorf("Room not found")
	}

	if err := room.SetGrid(client, gridData); err != nil {
		return err
	}

	return nil
}

func (e *Engine) Turn(room *db.Room, client *db.Client, x, y uint) (TurnResult, error) {
	return TurnResultMiss, nil
}
