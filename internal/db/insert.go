package db

import (
	"fmt"
)

func (m *manager) Insert(link string) error {
	m.Lock()
	defer func() {
		m.Unlock()
	}()

	m.close()
	defer func() {
		m.close()
	}()
	if err := m.open(); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	_, err := m.db.Exec("insert into links (link, status) values ($1, $2)", link, StatusNEW)
	if err != nil {
		return fmt.Errorf("db.Insert(exec): %w", err)
	}

	return nil
}
