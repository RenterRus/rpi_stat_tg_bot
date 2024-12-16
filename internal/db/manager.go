package db

import (
	"database/sql"
	"fmt"

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
	SelectAll(where string) ([]string, error)
	Update(link, status string) error
	Delete() error
}

type links struct {
	link string
}

type manager struct {
	pathToDB string
	db       *sql.DB
	block    chan struct{}
}

func NewManager(pathToDB string) Queue {
	res := &manager{
		pathToDB: pathToDB,
		block:    make(chan struct{}, 1),
	}

	work, _ := res.SelectAll(StatusWORK)

	for _, v := range work {
		res.Update(v, StatusNEW)
	}

	return res
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
