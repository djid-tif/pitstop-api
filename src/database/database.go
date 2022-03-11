package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"pitstop-api/src/utils"
)

const databaseFile = "env/database.db"

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		utils.PrintError(err)
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA journal_mode=WAL;")
	return db, nil
}
