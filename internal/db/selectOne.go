package db

import "fmt"

func (m *manager) SelectOne() (string, error) {
	m.block <- struct{}{}
	defer func() {
		<-m.block
	}()

	defer func() {
		m.close()
	}()
	if err := m.open(); err != nil {
		return "", fmt.Errorf("db.SelectOne: %w", err)
	}

	rows, err := m.db.Query("select link from links where status = $1 limit 1", StatusNEW)
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
