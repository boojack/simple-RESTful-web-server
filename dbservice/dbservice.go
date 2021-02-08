package dbservice

import (
	"database/sql"
	"fmt"

	// pg driver
	_ "github.com/lib/pq"
)

// DB db client
var DB *sql.DB = nil

func init() {
	connStr := "postgres://postgres:root@localhost/flomo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("connect failed")
		panic(err.Error())
	} else {
		DB = db
		fmt.Println("connect to mysql succeed")
	}

	// NOTE: 项目运行时无需 close
	// defer DB.Close()
}
