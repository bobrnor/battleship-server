package core

import (
	"sync"

	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
	"go.uber.org/zap"
)

type Room struct {
	sync.RWMutex
	uid string
}

func (r *Room) SetGrid(client *db.Client, gridData [13]uint8) error {
	r.Lock()
	defer r.Unlock()

	dbEntry, err := db.FindRoomByUID(nil, r.uid)
	if err != nil {
		return err
	}

	if dbEntry.State != db.RoomStateInitial {
		return errors.Errorf("Set grid operation available only during initial state")
	}

	grid := db.Grid{
		RoomID:   dbEntry.ID,
		ClientID: client.ID,
		Grid:     gridData,
	}
	if err := grid.Save(nil); err != nil {
		return err
	}

	return r.updateRoomState(dbEntry)
}

func (r *Room) Opponent(client *db.Client) (*db.Client, error) {
	clients, err := db.FindClientsByRoomUID(r.uid)
	if err != nil {
		return nil, err
	}

	for _, c := range clients {
		if c.ID != client.ID {
			return &c, nil
		}
	}

	return nil, nil
}

func (r *Room) OpponentGrid(client *db.Client) (*Grid, error) {
	gridsEntry, err := db.FindGridsByRoom(nil, client.ID)
	if err != nil {
		return nil, err
	}

	for _, gridEntry := range gridsEntry {
		if gridEntry.ClientID != client.ID {
			return &Grid{
				Data: gridEntry.Grid,
			}, nil
		}
	}

	return nil, nil
}

func (r *Room) IsReady() bool {
	r.RLock()
	defer r.RUnlock()

	dbEntry, err := db.FindRoomByUID(nil, r.uid)
	if err != nil {
		zap.S().Error(err)
		return false
	}

	return dbEntry.State == db.RoomStateReady
}

func (r *Room) updateRoomState(dbEntry *db.Room) error {
	grids, err := db.FindGridsByRoom(nil, dbEntry.ID)
	if err != nil {
		return err
	}

	clients, err := dbEntry.Clients(nil)
	if err != nil {
		return err
	}

	if len(grids) == 2 {
		dbEntry.State = db.RoomStateReady
		if err := dbEntry.Save(nil); err != nil {
			return err
		}

		for idx, client := range clients {
			var action string
			if idx == 0 {
				action = "turn"
			} else {
				action = "wait"
			}
			longpoll.DefaultLongpoll().Send(client.UID, map[string]interface{}{
				"action": action,
			})
		}
	}

	return nil
}
