package game

import (
	"sync"

	"github.com/pkg/errors"

	"git.nulana.com/bobrnor/battleship-server/db"
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

	if dbEntry.State != db.InitialState {
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

func (r *Room) updateRoomState(dbEntry *db.Room) error {
	grids, err := db.FindGridsByRoom(nil, dbEntry.ID)
	if err != nil {
		return err
	}

	if len(grids) == 2 {
		dbEntry.State = db.ReadyState
		return dbEntry.Save(nil)
	}

	return nil
}
