package client

import (
	sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"
	"github.com/pkg/errors"
)

type Client struct {
	ID  int64  `column:"id"`
	UID string `column:"uid"`
}

var (
	insert    *sqlsugar.InsertQuery
	findByUID *sqlsugar.SelectQuery
	findByID  *sqlsugar.SelectQuery
)

func init() {
	insert = sqlsugar.Insert((*Client)(nil)).Into("clients")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	findByUID = sqlsugar.Select((*Client)(nil)).From([]string{"clients"}).Where("uid = ?")
	if findByUID.Error() != nil {
		panic(findByUID.Error())
	}

	findByID = sqlsugar.Select((*Client)(nil)).From([]string{"clients"}).Where("id = ?")
	if findByID.Error() != nil {
		panic(findByID.Error())
	}
}

func FindByID(id int64) (*Client, error) {
	i, err := findByID.QueryRow(nil, id)
	if err != nil {
		return nil, err
	}

	var client *Client
	if i != nil {
		client = i.(*Client)
	}
	return client, nil
}

func FindByUID(uid string) (*Client, error) {
	i, err := findByUID.QueryRow(nil, uid)
	if err != nil {
		return nil, err
	}

	var client *Client
	if i != nil {
		client = i.(*Client)
	}
	return client, nil
}

func (c *Client) Save() error {
	results, err := insert.Exec(nil, c)
	if err == nil {
		c.ID, err = results.LastInsertId()
		err = errors.Wrap(err, "Can't get last insert id")
	}
	return err
}
