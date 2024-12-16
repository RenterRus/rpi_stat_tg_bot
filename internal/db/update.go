package db

import "fmt"

func (m *manager) Update(link, status string) error {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}
	defer m.close()

	_, err := m.db.Exec("update links set status = $1 where link = $2", status, link)
	if err != nil {
		return fmt.Errorf("db.Update(exec): %w", err)
	}

	return nil
}
