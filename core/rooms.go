package core

import (
	"time"

	"go.uber.org/zap"

	"sync"

	"git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
)

type Rooms struct {
	sync.Mutex
	rooms map[string]*Room
}

const (
//dispatcherTick      = 5 * time.Second
//confirmationTimeout = 10 * time.Minute
)

var (
	WrongClientNumber = errors.New("Wrang client numbers")
)

func NewRooms() *Rooms {
	return &Rooms{
		rooms: map[string]*Room{},
	}
}

func (r *Rooms) Register(clients []db.Client) (string, error) {
	zap.S().Infof("Register %+v", clients)

	if len(clients) != 2 {
		return "", errors.WithStack(WrongClientNumber)
	}

	tx, err := sqlsugar.Begin()
	if err != nil {
		return "", err
	}
	defer sqlsugar.RollbackOnRecover(tx, func(err error) {
		zap.S().Errorf("Can't register room %+v", err.Error())
	})

	room := db.Room{
		UID:   uuid.TimeOrderedUUID(),
		State: db.RoomStateInitial,
		TS:    time.Now().UTC(),
	}

	if err := room.Save(tx); err != nil {
		return "", err
	}

	if err := room.SetClients(tx, clients); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return room.UID, nil
}

func (r *Rooms) Room(uid string) *Room {
	r.Lock()
	defer r.Unlock()

	room, ok := r.rooms[uid]
	if !ok {
		room = &Room{
			uid: uid,
		}
		r.rooms[uid] = room
	}

	return room
}

//func (r *Rooms) dispatcherLoop() {
//	for range time.Tick(dispatcherTick) {
//		zap.S().Info("Dispatcher tick")
//		if err := db.FailUnconfirmedRooms(nil, confirmationTimeout); err != nil {
//			zap.S().Errorf("Can't delete unconfirmed rooms %+v", err.Error())
//		}
//	}
//}
