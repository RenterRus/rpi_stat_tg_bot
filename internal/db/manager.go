package db

import (
	"database/sql"
	"fmt"
	"sync"
)

type Queue interface {
	Insert(link string) error
	SelectOne() (string, error)
	SelectAll() ([]string, error)
	Update(link, status string) error
	Delete(link string) error
}

type links struct {
	link string
}

type manager struct {
	sync.Mutex
	pathToDB string
	db       *sql.DB
}

func NewManager(pathToDB string) Queue {
	return &manager{
		pathToDB: pathToDB,
	}
}

func (m *manager) open() error {
	var err error
	m.db, err = sql.Open("sqlite3", m.pathToDB)
	if err != nil {
		return fmt.Errorf("db.open: %w", err)
	}

	return nil
}

func (m *manager) close() {
	m.db.Close()
}
