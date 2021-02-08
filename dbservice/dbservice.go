package dbservice

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DB db client
var DB *sql.DB = nil

func init() {
	db, err := sql.Open("mysql", "root:root@/mine_flomo")
	if err != nil {
		fmt.Println("connect failed")
		panic(err.Error())
	} else {
		DB = db
		fmt.Println("connect succeed")
	}
	fmt.Println("db service inited")
}
