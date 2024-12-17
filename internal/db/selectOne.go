package db

import "fmt"

func (m *manager) SelectOne() (string, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return "", fmt.Errorf("db.SelectOne: %w", err)
	}
	defer m.close()

	rows, err := m.db.Query("select link, name from links where status = $1 order by RANDOM() limit 1", StatusNEW)
	if err != nil {
		return "", fmt.Errorf("db.SelectOne(query): %w", err)
	}
	defer func() {
		rows.Close()
	}()

	rows.Next()
	p := links{}
	err = rows.Scan(&p.Link, &p.Name)
	if err != nil {
		return "", fmt.Errorf("db.SelectOne(Scan): %w", err)
	}

	return p.Link, nil

}
