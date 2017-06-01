package room

import (
	"database/sql"

	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Room struct {
	ID        int64  `column:"id"`
	UID       string `column:"uid"`
	ClientID1 int64  `column:"client_id1"`
	ClientID2 int64  `column:"client_id2"`
}

var (
	insert *sqlsugar.InsertQuery
	find   *sqlsugar.SelectQuery
)

func init() {
	insert = sqlsugar.Insert((*Room)(nil)).Into("rooms")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	find = sqlsugar.Select((*Room)(nil)).From([]string{"rooms"}).Where("uid = ?")
	if find.Error() != nil {
		panic(find.Error())
	}
}

func FindByUID(tx *sql.Tx, uid string) (*Room, error) {
	i, err := find.QueryRow(tx, uid)
	if err != nil {
		return nil, err
	}
	return i.(*Room), nil
}

func (r *Room) Save(tx *sql.Tx) error {
	results, err := insert.Exec(tx, r)
	if err == nil {
		r.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insert id")
	}
	return err
}
