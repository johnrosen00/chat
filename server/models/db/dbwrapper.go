package db

import (
	"database/sql"
)

type Connection struct {
	db *sql.DB
}

func InitConnection(db *sql.DB) *Connection {
	return &Connection{db: db}
}
