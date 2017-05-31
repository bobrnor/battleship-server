package request

import sqlsugar "git.nulana.com/bobrnor/sqlsugar.git"

type Request struct {
	ID       int64 `column:"id"`
	ClientID int64 `column:"client_id"`
}

var (
	insert *sqlsugar.InsertQuery
	find   *sqlsugar.SelectQuery
)

func init() {
	insert = sqlsugar.Insert((*Request)(nil)).Into("requests")
	if insert.Error() != nil {
		panic(insert.Error())
	}

	find = sqlsugar.Select((*Request)(nil)).From([]string{"requests"})
	if find.Error() != nil {
		panic(find.Error())
	}
}

func All() ([]Request, error) {
	i, err := find.Query(nil)
	if err != nil {
		return []Request{}, err
	}
	return i.([]Request), nil
}

func (r *Request) Save() error {
	results, err := insert.Exec(nil, r)
	if err == nil {
		r.ID, err = results.LastInsertId()
	}
	return err
}
