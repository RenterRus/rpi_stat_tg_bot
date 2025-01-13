package db

import (
	"fmt"
	"strings"
)

func (m *manager) WorkList() string {
	queueList, err := m.SelectAll(StatusWORK)
	if err != nil {
		return fmt.Errorf("queue list: %w", err).Error()
	}

	workList, err := m.SelectAll(StatusNEW)
	if err != nil {
		return fmt.Errorf("work list: %w", err).Error()
	}

	queueList = append(queueList, workList...)

	var resp strings.Builder

	for _, v := range queueList {
		resp.WriteString(fmt.Sprintf("\"%s\",\n", v.Link))
	}

	return resp.String()
}
