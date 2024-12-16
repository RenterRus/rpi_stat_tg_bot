package db

import "fmt"

func (m *manager) Insert(link string) error {
	m.Lock()
	defer m.Unlock()

	if err := m.open(); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}
	defer m.close()

	result, err := m.db.Exec("insert into links (link, status) values ($1, 'NEW')", link)
	if err != nil {
		return fmt.Errorf("db.Insert(exec): %w", err)
	}

	fmt.Println(result.LastInsertId()) // id последнего добавленного объекта
	fmt.Println(result.RowsAffected()) // количество добавленных строк

	return nil
}
