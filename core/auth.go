package core

import (
	"database/sql"

	"git.nulana.com/bobrnor/battleship-server/db"
	"github.com/pkg/errors"
)

type auther struct {
	clientUID string
	client    *db.Client
	err       error
}

func AuthClient(uid string) (*db.Client, error) {
	a := auther{
		clientUID: uid,
	}
	return a.auth()
}

func (a *auther) auth() (*db.Client, error) {
	a.fetchClient()
	if a.client == nil {
		a.createClient()
	}
	return a.client, a.err
}

func (a *auther) fetchClient() {
	c, err := db.FindClientByUID(a.clientUID)
	if err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			a.err = err
		}
		return
	}
	a.client = c
}

func (a *auther) createClient() {
	if a.err != nil {
		return
	}

	newClient := db.Client{
		UID: a.clientUID,
	}

	if err := newClient.Save(nil); err != nil {
		a.err = err
		return
	}

	a.client = &newClient
}
