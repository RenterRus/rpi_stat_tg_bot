package db

import "fmt"

func (m *manager) Update(link, status string, name *string) error {
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

	if name != nil {
		_, err := m.db.Exec("update links set name = $1 where link = $2", *name, link)
		if err != nil {
			return fmt.Errorf("db.Update name(exec): %w", err)
		}
	}

	return nil
}
