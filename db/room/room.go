package room

import (
	"database/sql"
	"time"

	"git.nulana.com/bobrnor/battleship-server/db/client"
	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	InitialState = iota
	ConfirmedState
	FailedState
)

type Room struct {
	ID    int64     `column:"id"`
	UID   string    `column:"uid"`
	State int64     `column:"state"`
	TS    time.Time `column:"ts"`

	clients []client.Client
}

type RoomClient struct {
	RoomID    int64 `column:"room_id"`
	ClientID  int64 `column:"client_id"`
	Confirmed bool  `column:"confirmed"`
}

var (
	insert *sqlsugar.InsertQuery
	update *sqlsugar.UpdateQuery
	fail   *sqlsugar.UpdateQuery
	find   *sqlsugar.SelectQuery

	deleteClients *sqlsugar.DeleteQuery
	insertClients *sqlsugar.InsertQuery
	confirmClient *sqlsugar.UpdateQuery
)

func init() {
	insert = sqlsugar.Insert((*Room)(nil)).Into("rooms")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	update = sqlsugar.Update("rooms").SetAll((*Room)(nil)).Where("id = ?")
	if update.Error() != nil {
		panic(update.Error())
	}

	fail = sqlsugar.UpdateMultiple([]string{"rooms", "room_clients", "clients"}).Set([]string{"rooms.state"}).Where("room_clients.room_id = rooms.id AND room_clients.confirmed = 0 AND DATE_ADD(rooms.ts, INTERVAL ? SECOND) < UTC_TIMESTAMP()")
	if fail.Error() != nil {
		panic(fail.Error())
	}

	find = sqlsugar.Select((*Room)(nil)).From([]string{"rooms"}).Where("uid = ?")
	if find.Error() != nil {
		panic(find.Error())
	}

	deleteClients = sqlsugar.Delete("room_clients").Where("room_id = ?")
	if deleteClients.Error() != nil {
		panic(deleteClients.Error())
	}

	insertClients = sqlsugar.Insert((*RoomClient)(nil)).Into("room_clients")
	if insertClients.Error() != nil {
		panic(insertClients.Error())
	}

	confirmClient = sqlsugar.Update("room_clients").Set([]string{"confirmed"}).Where("room_id = ? AND client_id = ?")
	if confirmClient.Error() != nil {
		panic(confirmClient.Error())
	}
}

func FindByUID(tx *sql.Tx, uid string) (*Room, error) {
	i, err := find.QueryRow(tx, uid)
	if err != nil {
		return nil, err
	}
	if i != nil {
		return i.(*Room), nil
	}
	return nil, nil
}

func FailUnconfirmed(tx *sql.Tx, timeout time.Duration) error {
	room := &Room{
		State: FailedState,
	}
	_, err := fail.Exec(tx, room, timeout/time.Second)
	return err
}

func (r *Room) Clients(tx *sql.Tx) ([]client.Client, error) {
	if len(r.clients) == 0 {
		clients, err := client.FindByRoomID(r.ID)
		if err != nil {
			return r.clients, err
		}
		r.clients = clients
	}
	return r.clients, nil
}

func (r *Room) SetClients(tx *sql.Tx, clients []client.Client) error {
	localTx := false
	if tx == nil {
		newTx, err := sqlsugar.Begin()
		if err != nil {
			return err
		}

		localTx = true
		tx = newTx
		defer sqlsugar.RollbackOnRecover(tx, func(err error) {
			zap.S().Errorf("Error while setting room clients %+v", err.Error())
		})
	}

	if _, err := deleteClients.Exec(tx, r.ID); err != nil {
		tx.Rollback()
		return err
	}

	// TODO: do it better (for 1 query)
	for _, client := range clients {
		roomClient := &RoomClient{
			RoomID:    r.ID,
			ClientID:  client.ID,
			Confirmed: false,
		}
		if _, err := insertClients.Exec(tx, roomClient); err != nil {
			tx.Rollback()
			return err
		}
	}

	if localTx {
		tx.Commit()
	}

	r.clients = clients

	return nil
}

func (r *Room) Confirm(tx *sql.Tx, c *client.Client) error {
	roomClient := RoomClient{
		Confirmed: true,
	}
	_, err := confirmClient.Exec(nil, &roomClient, r.ID, c.ID)
	return err
}

func (r *Room) Save(tx *sql.Tx) error {
	if r.ID > 0 {
		return r.update(tx)
	} else {
		return r.insert(tx)
	}
}

func (r *Room) insert(tx *sql.Tx) error {
	results, err := insert.Exec(tx, r)
	if err == nil {
		r.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insert id")
	}
	return err
}

func (r *Room) update(tx *sql.Tx) error {
	_, err := update.Exec(tx, r, r.ID)
	return err
}
