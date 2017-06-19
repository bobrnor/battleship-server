package db

import (
	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Session struct {
	ID        int64  `column:"id"`
	ClientID  string `column:"client_id"`
	SessionID string `column:"session_id"`
}

var (
	insertSession *sqlsugar.InsertQuery
	updateSession *sqlsugar.UpdateQuery
)

func init() {
	insertSession = sqlsugar.Insert((*Session)(nil)).Into("sessions")
	if insertSession.Error() != nil {
		panic(insertSession.Error())
	}

	updateSession = sqlsugar.Update("sessions").SetAll((*Session)(nil)).Where("`id` = ?")
	if updateSession.Error() != nil {
		panic(updateSession.Error())
	}
}

func (s *Session) Save() error {
	if s.ID > 0 {
		_, err := updateSession.Exec(nil, s, s.ID)
		return err
	} else {
		results, err := insertSession.Exec(nil, s)
		if err == nil {
			s.ID, err = results.LastInsertId()
			err = errors.Wrap(err, "Can't get last insertSession id")
		}
		return err
	}
}
