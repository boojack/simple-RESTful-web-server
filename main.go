package main

import (
	"log"
	"neosmemo/backend/router"
	"net/http"
)

func init() {
	// do nth
}

func main() {
	router := router.Router

	log.Fatal(http.ListenAndServe(":8080", router))
}
