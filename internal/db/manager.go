package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	StatusNEW  = "NEW"
	StatusWORK = "WORK"
	StatusDONE = "DONE"
)

type Queue interface {
	Insert(link string) error
	SelectOne() (string, error)
	SelectAll(where string) ([]links, error)
	Update(link, status string, name *string) error
	Delete() error
}

type links struct {
	Link string
	Name *string
}

type manager struct {
	pathToDB string
	db       *sql.DB
	sync.Mutex
}

func NewManager(pathToDB string) Queue {
	res := &manager{
		pathToDB: pathToDB,
	}

	work, _ := res.SelectAll(StatusWORK)

	for _, v := range work {
		res.Update(v.Link, StatusNEW, nil)
	}

	return res
}

func (m *manager) open() error {
	var err error
	m.db, err = sql.Open("sqlite3", m.pathToDB)
	if err != nil {
		return fmt.Errorf("db.open: %w", err)
	}

	err = m.db.Ping()
	if err != nil {
		return fmt.Errorf("db.open(ping): %w", err)
	}

	return nil
}

func (m *manager) close() {
	m.db.Close()
}
