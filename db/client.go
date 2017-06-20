package db

import (
	"database/sql"

	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Client struct {
	ID  int64  `column:"id"`
	UID string `column:"uid"`
}

var (
	insertClient        *sqlsugar.InsertQuery
	findClientByUID     *sqlsugar.SelectQuery
	findClientByID      *sqlsugar.SelectQuery
	findClientByRoomID  *sqlsugar.SelectQuery
	findClientByRoomUID *sqlsugar.SelectQuery
)

func init() {
	insertClient = sqlsugar.Insert((*Client)(nil)).Into("clients")
	if insertClient.Error() != nil {
		panic(insertClient.Error())
	}

	findClientByUID = sqlsugar.Select((*Client)(nil)).From([]string{"clients"}).Where("uid = ?")
	if findClientByUID.Error() != nil {
		panic(findClientByUID.Error())
	}

	findClientByID = sqlsugar.Select((*Client)(nil)).From([]string{"clients"}).Where("id = ?")
	if findClientByID.Error() != nil {
		panic(findClientByID.Error())
	}

	findClientByRoomID = sqlsugar.Select((*Client)(nil)).From([]string{"clients", "room_clients"}).Where("room_clients.room_id = ? && clients.id = room_clients.client_id")
	if findClientByRoomID.Error() != nil {
		panic(findClientByRoomID.Error())
	}

	findClientByRoomUID = sqlsugar.Select((*Client)(nil)).From([]string{"clients", "rooms", "room_clients"}).Where("rooms.uid = ? AND room_clients.room_id = rooms.id && clients.id = room_clients.client_id")
	if findClientByRoomUID.Error() != nil {
		panic(findClientByRoomUID.Error())
	}
}

func FindClientByID(id int64) (*Client, error) {
	i, err := findClientByID.QueryRow(nil, id)
	if err != nil {
		return nil, err
	}

	var client *Client
	if i != nil {
		client = i.(*Client)
	}
	return client, nil
}

func FindClientByUID(uid string) (*Client, error) {
	i, err := findClientByUID.QueryRow(nil, uid)
	if err != nil {
		return nil, err
	}

	var client *Client
	if i != nil {
		client = i.(*Client)
	}
	return client, nil
}

func FindClientByRoomID(roomID int64) ([]Client, error) {
	i, err := findClientByRoomID.Query(nil, roomID)
	if err != nil {
		return nil, err
	}

	var c []Client
	if i != nil {
		c = i.([]Client)
	}
	return c, nil
}

func FindClientByRoomUID(roomUID string) ([]Client, error) {
	i, err := findClientByRoomUID.Query(nil, roomUID)
	if err != nil {
		return nil, err
	}

	var c []Client
	if i != nil {
		c = i.([]Client)
	}
	return c, nil
}

func (c *Client) Save(tx *sql.Tx) error {
	results, err := insertClient.Exec(tx, c)
	if err == nil {
		c.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insertClient id")
	}
	return err
}
