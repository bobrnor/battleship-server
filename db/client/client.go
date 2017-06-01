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
	insert *sqlsugar.InsertQuery
	find   *sqlsugar.SelectQuery
)

func init() {
	insert = sqlsugar.Insert((*Client)(nil)).Into("clients")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	find = sqlsugar.Select((*Client)(nil)).From([]string{"clients"}).Where("uid = ?")
	if find.Error() != nil {
		panic(find.Error())
	}
}

func FindByUID(uid string) (*Client, error) {
	i, err := find.QueryRow(nil, uid)
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
