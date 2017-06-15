package rooms

import (
	"time"

	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"

	"git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
)

type Rooms struct{}

const (
	dispatcherTick      = 5 * time.Second
	confirmationTimeout = 10 * time.Minute
)

var (
	rooms *Rooms

	WrongClientNumber = errors.New("Wrang client numbers")
)

func init() {
	rooms = &Rooms{}
	go rooms.dispatcherLoop()
}

func DefaultRooms() *Rooms {
	return rooms
}

func (r *Rooms) Register(clients []client.Client) (string, error) {
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

	room := room.Room{
		UID:   uuid.TimeOrderedUUID(),
		State: room.InitialState,
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

func (r *Rooms) dispatcherLoop() {
	for range time.Tick(dispatcherTick) {
		zap.S().Info("Dispatcher tick")
		if err := room.FailUnconfirmed(nil, confirmationTimeout); err != nil {
			zap.S().Errorf("Can't delete unconfirmed rooms %+v", err.Error())
		}
	}
}
