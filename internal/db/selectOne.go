package db

import "fmt"

func (m *manager) SelectOne() (string, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return "", fmt.Errorf("db.SelectOne: %w", err)
	}
	defer m.close()

	rows, err := m.db.Query("select link from links where status = 'NEW' limit 1")
	if err != nil {
		return "", fmt.Errorf("db.SelectOne(query): %w", err)
	}

	rows.Next()
	p := links{}
	err = rows.Scan(&p.link)
	if err != nil {
		return "", fmt.Errorf("db.SelectOne(Scan): %w", err)
	}

	return p.link, nil

}
