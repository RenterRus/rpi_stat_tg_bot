package db

import "fmt"

// Очищаем ссылки из таблицы, которые уже были скачаны. Они уже не играют функциональной роли
func (m *manager) Delete() error {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return fmt.Errorf("db.Delete: %w", err)
	}
	defer m.close()

	_, err := m.db.Exec("delete from links where status = $1", StatusDONE)
	if err != nil {
		return fmt.Errorf("db.Delete(exec): %w", err)
	}

	return nil
}

func (m *manager) DeleteByLink(link string) error {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return fmt.Errorf("db.Delete: %w", err)
	}
	defer m.close()

	_, err := m.db.Exec("delete from links where link = $1", link)
	if err != nil {
		return fmt.Errorf("db.Delete(exec): %w", err)
	}

	return nil
}
