package core

import (
	"git.nulana.com/bobrnor/battleship-grid.git"
	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
	"github.com/pkg/errors"
)

type Engine struct{}

type TurnResult uint

const (
	TurnResultMiss = TurnResult(iota)
	TurnResultHit
	TurnResultLose
	TurnResultWin
)

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) SetGrid(dbRoom *db.Room, client *db.Client, gridData [13]byte) error {
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

	opponent, err := db.FindClientByRoomIDAndNotClientID(dbRoom.ID, client.ID)
	if err != nil {
		return TurnResultMiss, err
	}

	opponentGrid, err := db.FindGridByRoomAndClient(nil, dbRoom.ID, opponent.ID)
	if err != nil {
		return TurnResultMiss, err
	}

	if err := e.storeHit(opponentGrid, x, y); err != nil {
		return TurnResultMiss, err
	}

	result := TurnResultMiss
	if (&opponentGrid.Grid).Get(x, y) {
		result = TurnResultHit
	}

	if result == TurnResultHit && e.isEnded(opponentGrid) {
		result = TurnResultWin
		longpoll.DefaultLongpoll().Send(opponent.UID, map[string]interface{}{
			"type":   "game_over",
			"action": "lose",
		})
	} else {
		longpoll.DefaultLongpoll().Send(opponent.UID, map[string]interface{}{
			"type": "opponent_turn",
			"x":    x,
			"y":    y,
		})
	}

	return result, nil
}

func (e *Engine) storeHit(g *db.Grid, x, y uint) error {
	(&g.Hits).Set(x, y)
	return g.Save(nil)
}

func (e *Engine) isEnded(g *db.Grid) bool {
	diff := grid.Diff(&g.Grid, &g.Hits)
	return diff.IsEmpty()
}
