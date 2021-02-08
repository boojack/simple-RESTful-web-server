package main

import (
	"log"
	"neosmemo/backend/dbservice"
	"neosmemo/backend/router"
	"net/http"
)

func init() {
	if router.Router == nil {
		panic("router init failed")
	}
	if dbservice.DB == nil {
		panic("dbservice init failed")
	}
}

func main() {
	r := router.Router
	r.ServeFiles("/web/*filepath", http.Dir("./static/"))

	log.Fatal(http.ListenAndServe(":8080", r))
}
