package db

import (
	"fmt"
)

func (m *manager) SelectAll(whereStatus string) ([]links, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return nil, fmt.Errorf("db.SelectAll: %w", err)
	}
	defer m.close()

	rows, err := m.db.Query("select link from links where status = $1", whereStatus)
	defer func() {
		rows.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("db.SelectAll(Query(where)): %w", err)
	}
	res := make([]links, 0, 2)
	for rows.Next() {
		p := links{}
		err = rows.Scan(&p.Link, &p.Name)
		if err != nil {
			break
		}

		res = append(res, p)
	}

	return res, nil
}
