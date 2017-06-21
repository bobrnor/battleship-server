package core

import (
	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
	"github.com/pkg/errors"
)

type Engine struct{}

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

func (e *Engine) Turn(dbRoom *db.Room, client *db.Client, x, y uint) (TurnResult, error) {
	rooms := MainRooms()
	room := rooms.Room(dbRoom.UID)
	if room == nil {
		return TurnResultMiss, errors.Errorf("Room not found")
	}

	if !room.IsReady() {
		return TurnResultMiss, errors.Errorf("Room is not ready")
	}

	opponentGrid, err := room.OpponentGrid(client)
	if err != nil {
		return TurnResultMiss, err
	}

	opponent, err := room.Opponent(client)
	if err != nil {
		return TurnResultMiss, err
	}

	longpoll.DefaultLongpoll().Send(opponent.UID, map[string]interface{}{
		"x": x,
		"y": y,
	})

	return opponentGrid.Turn(x, y), nil
}
