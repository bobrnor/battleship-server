package core

import (
	"sync"

	"github.com/pkg/errors"

	"log"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/longpoll.git"
)

type Room struct {
	sync.RWMutex
	uid string
}

func (r *Room) SetGrid(client *db.Client, gridData [13]byte) error {
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
		Hits:     [13]byte{},
	}
	if err := grid.Save(nil); err != nil {
		return err
	}

	return r.updateRoomState(dbEntry)
}

func (r *Room) IsReady() bool {
	r.RLock()
	defer r.RUnlock()

	dbEntry, err := db.FindRoomByUID(nil, r.uid)
	if err != nil {
		log.Printf("%+v", err.Error())
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
				"type":   "game",
				"action": action,
			})
		}
	}

	return nil
}
