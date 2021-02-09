package helper

import (
	"database/sql"
	"fmt"

	// pg driver
	_ "github.com/lib/pq"
)

// DBService db client
var DBService *sql.DB = nil

func init() {
	connStr := "postgres://postgres:root@localhost/flomo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("connect failed")
		panic(err.Error())
	} else {
		DBService = db
		fmt.Println("connect to postgres succeed")
	}

	// NOTE: As the official notes:
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	// defer DBService.Close()
}
