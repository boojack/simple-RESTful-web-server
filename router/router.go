package router

import (
	"neosmemo/backend/handler"
	"neosmemo/backend/handler/memo"
	"neosmemo/backend/handler/user"

	"github.com/julienschmidt/httprouter"
)

// Router com
var Router *httprouter.Router = nil

// NOTE: 在这里注册路由
func init() {
	Router = httprouter.New()

	// NOTE: 中间件的实现方式
	Router.GET("/hello/:id", handler.Middleware(user.GetUserByID, handler.MiddlewareConfig{Cors: true, JSON: true}))

	// user about
	Router.GET("/api/user/all", user.GetAllUser)
	Router.POST("/api/user/check", user.CheckUsernameUsed)
	Router.POST("/api/user/info", user.GetUserByID)
	Router.POST("/api/user/signup", user.DoSignUp)
	Router.POST("/api/user/signin", user.DoSignIn)

	// memo about
	Router.GET("/api/memo/", memo.GetAllMemos)
	Router.GET("/api/memo/:id", memo.GetMemoByID)
	Router.POST("/api/memo/new", memo.CreateMemo)
	Router.POST("/api/memo/update", memo.UpdateMemo)
	Router.POST("/api/memo/delete", memo.DeleteMemo)
}
