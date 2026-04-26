package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open db connection: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db: ", err)
	}

	fmt.Println("Connected to database!")
	return db
}
