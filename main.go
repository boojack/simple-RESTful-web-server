package main

import (
	"log"
	"neosmemo/backend/helper"
	"neosmemo/backend/router"
	"net/http"
)

func init() {
	if router.Router == nil {
		panic("router init failed")
	}
	if helper.DBService == nil {
		panic("dbservice init failed")
	}
}

func main() {
	r := router.Router

	// NOTE: static file server, use for test
	r.ServeFiles("/web/*filepath", http.Dir("./static/"))

	log.Fatal(http.ListenAndServe(":8080", r))
}
