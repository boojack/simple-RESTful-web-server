package router

import (
	"neosmemo/backend/controller/user"

	"github.com/julienschmidt/httprouter"
)

// Router com
var Router *httprouter.Router = nil

// NOTE: 在这里注册路由
func init() {
	Router = httprouter.New()

	Router.GET("/hello/:id", user.GetUserByID)
	Router.GET("/api/user/all", user.GetAllUser)
	Router.POST("/api/user/check", user.CheckUsernameUsed)
	Router.POST("/api/user/signup", user.SignUp)
}
