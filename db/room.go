package db

import (
	"database/sql"
	"time"

	"log"

	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Room struct {
	ID    int64     `column:"id"`
	UID   string    `column:"uid"`
	State int64     `column:"state"`
	TS    time.Time `column:"ts"`

	clients []Client
}

type RoomClient struct {
	RoomID   int64 `column:"room_id"`
	ClientID int64 `column:"client_id"`
}

const (
	RoomStateInitial = iota
	RoomStateReady
)

var (
	insertRoom *sqlsugar.InsertQuery
	updateRoom *sqlsugar.UpdateQuery
	failRoom   *sqlsugar.UpdateQuery
	findRoom   *sqlsugar.SelectQuery

	deleteClientsFromRoom *sqlsugar.DeleteQuery
	insertClientsToRoom   *sqlsugar.InsertQuery
)

func init() {
	insertRoom = sqlsugar.Insert((*Room)(nil)).Into("rooms")
	if insertRoom.Error() != nil {
		panic(insertRoom.Error())
	}

	updateRoom = sqlsugar.Update("rooms").SetAll((*Room)(nil)).Where("id = ?")
	if updateRoom.Error() != nil {
		panic(updateRoom.Error())
	}

	failRoom = sqlsugar.UpdateMultiple([]string{"rooms", "room_clients", "clients"}).Set([]string{"rooms.state"}).Where("room_clients.room_id = rooms.id AND room_clients.confirmed = 0 AND DATE_ADD(rooms.ts, INTERVAL ? SECOND) < UTC_TIMESTAMP()")
	if failRoom.Error() != nil {
		panic(failRoom.Error())
	}

	findRoom = sqlsugar.Select((*Room)(nil)).From([]string{"rooms"}).Where("uid = ?")
	if findRoom.Error() != nil {
		panic(findRoom.Error())
	}

	deleteClientsFromRoom = sqlsugar.Delete("room_clients").Where("room_id = ?")
	if deleteClientsFromRoom.Error() != nil {
		panic(deleteClientsFromRoom.Error())
	}

	insertClientsToRoom = sqlsugar.Insert((*RoomClient)(nil)).Into("room_clients")
	if insertClientsToRoom.Error() != nil {
		panic(insertClientsToRoom.Error())
	}
}

func FindRoomByUID(tx *sql.Tx, uid string) (*Room, error) {
	i, err := findRoom.QueryRow(tx, uid)
	if err != nil {
		return nil, err
	}
	return i.(*Room), nil
}

func (r *Room) Clients(tx *sql.Tx) ([]Client, error) {
	if len(r.clients) == 0 {
		clients, err := FindClientsByRoomID(r.ID)
		if err != nil {
			return r.clients, err
		}
		r.clients = clients
	}
	return r.clients, nil
}

func (r *Room) SetClients(tx *sql.Tx, clients []Client) error {
	localTx := false
	if tx == nil {
		newTx, err := sqlsugar.Begin()
		if err != nil {
			return err
		}

		localTx = true
		tx = newTx
		defer sqlsugar.RollbackOnRecover(tx, func(err error) {
			log.Printf("Error while setting room clients %+v", err.Error())
		})
	}

	if _, err := deleteClientsFromRoom.Exec(tx, r.ID); err != nil {
		tx.Rollback()
		return err
	}

	// TODO: do it better (for 1 query)
	for _, client := range clients {
		roomClient := &RoomClient{
			RoomID:   r.ID,
			ClientID: client.ID,
		}
		if _, err := insertClientsToRoom.Exec(tx, roomClient); err != nil {
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

func (r *Room) Save(tx *sql.Tx) error {
	if r.ID > 0 {
		return r.updateRoom(tx)
	} else {
		return r.insertRoom(tx)
	}
}

func (r *Room) insertRoom(tx *sql.Tx) error {
	results, err := insertRoom.Exec(tx, r)
	if err == nil {
		r.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insertRoom id")
	}
	return err
}

func (r *Room) updateRoom(tx *sql.Tx) error {
	_, err := updateRoom.Exec(tx, r, r.ID)
	return err
}
