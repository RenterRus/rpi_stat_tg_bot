package db

import "fmt"

// Очищаем ссылки из таблицы, которые уже были скачаны. Они уже не играют функциональной роли
func (m *manager) Delete() error {
	m.block <- struct{}{}
	defer func() {
		<-m.block
	}()

	defer func() {
		m.close()
	}()
	if err := m.open(); err != nil {
		return fmt.Errorf("db.Delete: %w", err)
	}

	_, err := m.db.Exec("delete from links where status = $1", StatusDONE)
	if err != nil {
		return fmt.Errorf("db.Delete(exec): %w", err)
	}

	return nil
}
