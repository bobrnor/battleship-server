package client

import sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"

type Client struct {
	ID  int64  `column:"id"`
	UID string `column:"uid"`
}

var (
	insert *sqlsugar.InsertQuery
	update *sqlsugar.UpdateQuery
	find   *sqlsugar.SelectQuery
)

func init() {
	insert = sqlsugar.Insert((*Client)(nil)).Into("clients")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	update := sqlsugar.Update("clients").SetAll((*Client)(nil)).Where("id = ?")
	if update.Error() != nil {
		panic(update.Error())
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
	if c.ID > 0 {
		_, err := update.Exec(nil, c, c.ID)
		return err
	} else {
		results, err := insert.Exec(nil, c)
		if err == nil {
			c.ID, err = results.LastInsertId()
		}
		return err
	}
}
