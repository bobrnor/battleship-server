package rooms

import (
	"time"

	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	"git.nulana.com/bobrnor/battleship-server/db/room"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/pkg/errors"
)

type Rooms struct{}

const (
	dispatcherTick      = 5 * time.Second
	confirmationTimeout = 1 * time.Minute
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

func (r *Rooms) Register(clients []*client.Client) (string, error) {
	zap.S().Infof("Register %+v", clients)

	if len(clients) != 2 {
		return "", errors.WithStack(WrongClientNumber)
	}

	// TODO: do it using N:N rel.
	room := room.Room{
		UID:          uuid.TimeOrderedUUID(),
		ClientID1:    clients[0].ID,
		ClientID2:    clients[1].ID,
		Client1State: room.InitialState,
		Client2State: room.InitialState,
		TS:           time.Now().UTC(),
	}

	if err := room.Save(nil); err != nil {
		return "", err
	}

	return room.UID, nil
}

func (r *Rooms) dispatcherLoop() {
	for range time.Tick(dispatcherTick) {
		zap.S().Info("Dispatcher tick")
		if err := room.DeleteUnconfirmed(nil, confirmationTimeout); err != nil {
			zap.S().Errorf("Can't delete unconfirmed rooms %+v", err.Error())
		}
	}
}
