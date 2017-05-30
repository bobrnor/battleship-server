package db

import (
	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"go.uber.org/zap"
)

type Session struct {
	ID        int64  `column:"id"`
	ClientID  string `column:"client_id"`
	SessionID string `column:"session_id"`
}

var (
	insert *sqlsugar.InsertQuery
	update *sqlsugar.UpdateQuery
)

func init() {
	insert = sqlsugar.Insert((*Session)(nil)).Into("sessions")
	if insert.Error() != nil {
		zap.S().Fatalw("Can't make select query", insert.Error())
	}

	update = sqlsugar.Update("sessions").SetAll((*Session)(nil)).Where("`id` = ?")
	if update.Error() != nil {
		zap.S().Fatalw("Can't make select query", update.Error())
	}
}

func (s *Session) Save() error {
	if s.ID > 0 {
		_, err := update.Exec(nil, s, s.ID)
		return err
	} else {
		results, err := insert.Exec(nil, s)
		if err == nil {
			s.ID, err = results.LastInsertId()
		}
		return err
	}
}
