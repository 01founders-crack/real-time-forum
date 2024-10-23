package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var MyDB *sql.DB

func Init() {
	var err error
	MyDB, err = sql.Open("sqlite3", "server/database/forum.db")
	if err != nil {
		log.Fatal("Invalid DB config, unable to open database:", err)
	}
}
