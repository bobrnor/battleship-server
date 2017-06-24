package db

import (
	"database/sql"

	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Grid struct {
	ID       int64  `column:"id"`
	RoomID   int64  `column:"room_id"`
	ClientID int64  `column:"client_id"`
	Grid     []byte `column:"grid"`
}

var (
	insertGrid       *sqlsugar.InsertQuery
	updateGrid       *sqlsugar.UpdateQuery
	findGrid         *sqlsugar.SelectQuery
	findOpponentGrid *sqlsugar.SelectQuery
	findAllGrid      *sqlsugar.SelectQuery
)

func init() {
	insertGrid = sqlsugar.Insert((*Grid)(nil)).Into("grids")
	if insertGrid.Error() != nil {
		panic(insertGrid.Error())
	}

	updateGrid = sqlsugar.Update("grid").Set([]string{"grid"}).Where("room_id = ? AND client_id = ?")
	if updateGrid.Error() != nil {
		panic(updateGrid.Error())
	}

	findGrid = sqlsugar.Select((*Grid)(nil)).From([]string{"grids"}).Where("room_id = ? AND client_id = ?")
	if findGrid.Error() != nil {
		panic(findGrid.Error())
	}

	findOpponentGrid = sqlsugar.Select((*Grid)(nil)).From([]string{"grids"}).Where("room_id = ? AND client_id != ?")
	if findOpponentGrid.Error() != nil {
		panic(findOpponentGrid.Error())
	}

	findAllGrid = sqlsugar.Select((*Grid)(nil)).From([]string{"grids"}).Where("room_id = ?")
	if findGrid.Error() != nil {
		panic(findGrid.Error())
	}
}

func FindGridByRoomAndClient(tx *sql.Tx, roomID, clientID int64) (*Grid, error) {
	i, err := findGrid.QueryRow(tx, roomID, clientID)
	if err != nil {
		return nil, err
	}
	return i.(*Grid), nil
}

func FindGridByRoomAndNotClient(tx *sql.Tx, roomID, clientID int64) (*Grid, error) {
	i, err := findOpponentGrid.QueryRow(tx, roomID, clientID)
	if err != nil {
		return nil, err
	}
	return i.(*Grid), nil
}

func FindGridsByRoom(tx *sql.Tx, roomID int64) ([]Grid, error) {
	i, err := findAllGrid.Query(tx, roomID)
	if err != nil {
		return nil, err
	}
	return i.([]Grid), nil
}

func (g *Grid) Save(tx *sql.Tx) error {
	if g.ID > 0 {
		return g.updateGrid(tx)
	} else {
		return g.insertGrid(tx)
	}
}

func (g *Grid) insertGrid(tx *sql.Tx) error {
	results, err := insertGrid.Exec(tx, g)
	if err == nil {
		g.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insertGrid id")
	}
	return err
}

func (g *Grid) updateGrid(tx *sql.Tx) error {
	_, err := updateGrid.Exec(tx, g, g.ID)
	return err
}
