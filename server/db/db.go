package db

import (
	"database/sql"
	_ "embed"
	"os"

	"cloud-server/logger"

	_ "github.com/mattn/go-sqlite3"
)

func IsTest() bool {
	return os.Getenv("TEST") == "true"
}

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
	log.Info("Starting handler")
	var err error
	dataSource := "storage/storage.db"

	if IsTest() {
		log.Warn("Running in test mode")
		dataSource = "storage/storage_test.db"
	}

	db.Connection, err = sql.Open("sqlite3", dataSource)
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
	log.Warn("Handler stopped")
	if IsTest() {
		log.Warn("Removing test database")
		os.Remove("storage/storage_test.db")
	}
	log.Close()
}
