package db

import (
	"fmt"
	"strings"
)

func (m *manager) Update(link, status string, name *string) error {
	m.Lock()
	defer m.Unlock()

	const MAX_EXT_SIZE = 4

	if err := m.open(); err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}
	defer m.close()

	_, err := m.db.Exec("update links set status = $1 where link = $2", status, link)
	if err != nil {
		return fmt.Errorf("db.Update(exec): %w", err)
	}

	if name != nil {
		finalName := ""
		for i, v := range strings.Split(*name, ".") {
			if i > 1 {
				if len(v) > MAX_EXT_SIZE {
					finalName += fmt.Sprintf(". %s", v)
				}
			} else {
				finalName += v
			}
		}
		_, err := m.db.Exec("update links set name = $1 where link = $2", finalName, link)
		if err != nil {
			return fmt.Errorf("db.Update name(exec): %w", err)
		}
	}

	return nil
}
