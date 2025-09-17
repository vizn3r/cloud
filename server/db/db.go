package db

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Connection *sql.DB
}

func NewDB() *DB {
	return &DB{nil}
}

//go:embed sql/tables.sql
var tableQuery string

func (db *DB) Start() {
	log.Println("Starting DB handler")
	var err error

	dbPath := filepath.Join("storage", "storage.db")
	db.Connection, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating tables")
	log.Println(tableQuery)
	_, err = db.Connection.Exec(tableQuery)
	if err != nil {
		db.Connection.Close()
		log.Fatal(err)
	}
}

func (db *DB) Stop() {
	db.Connection.Close()
	log.Println("DB handler stopped")
}
