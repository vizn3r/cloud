package db

import (
	"cloud-server/logger"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var log = logger.New(" DB ", logger.Magenta)

type DB struct {
	Connection *sql.DB
}

func NewDB() *DB {
	return &DB{nil}
}

//go:embed sql/tables.sql
var tableQuery string

func (db *DB) Start() {
	log.Info("Starting DB handler")
	var err error

	db.Connection, err = sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Creating tables...")
	_, err = db.Connection.Exec(tableQuery)
	if err != nil {
		db.Connection.Close()
		log.Fatal(err)
	}
}

func (db *DB) Stop() {
	db.Connection.Close()
	log.Warn("DB handler stopped")
	log.Close()
}
