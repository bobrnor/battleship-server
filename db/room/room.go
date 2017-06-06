package room

import (
	"database/sql"
	"time"

	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

const (
	InitialState = iota
	ConfirmedState
)

type Room struct {
	ID           int64     `column:"id"`
	UID          string    `column:"uid"`
	ClientID1    int64     `column:"client_id1"`
	ClientID2    int64     `column:"client_id2"`
	Client1State int64     `column:"client1_state"`
	Client2State int64     `column:"client2_state"`
	TS           time.Time `column:"ts"`
}

var (
	insert *sqlsugar.InsertQuery
	update *sqlsugar.UpdateQuery
	find   *sqlsugar.SelectQuery
	delete *sqlsugar.DeleteQuery
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

	find = sqlsugar.Select((*Room)(nil)).From([]string{"rooms"}).Where("uid = ?")
	if find.Error() != nil {
		panic(find.Error())
	}

	delete = sqlsugar.Delete("rooms").Where("(client1_state = ? || client2_state = ?) AND DATE_ADD(ts, INTERVAL ? SECOND) < UTC_TIMESTAMP()")
	if delete.Error() != nil {
		panic(delete.Error())
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

func DeleteUnconfirmed(tx *sql.Tx, timeout time.Duration) error {
	_, err := delete.Exec(tx, InitialState, InitialState, timeout/time.Second)
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
