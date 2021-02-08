package router

import (
	"neosmemo/backend/controller/user"

	"github.com/julienschmidt/httprouter"
)

// Router com
var Router *httprouter.Router = nil

func init() {
	Router = httprouter.New()

	Router.GET("/hello/:id", user.GetUserByID)
}
