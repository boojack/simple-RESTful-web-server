package router

import (
	"neosmemo/backend/handler"
	"neosmemo/backend/handler/memo"
	"neosmemo/backend/handler/user"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router router
var Router *httprouter.Router = nil

// NOTE: 在这里注册路由
func init() {
	Router = httprouter.New()

	// user about
	Router.GET("/api/user/me", handler.Middleware(user.GetMyUserInfo, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/user/signup", handler.Middleware(user.DoSignUp, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/user/signin", handler.Middleware(user.DoSignIn, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/user/signout", handler.Middleware(user.DoSignOut, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/user/check", handler.Middleware(user.CheckUsernameUsed, handler.MiddlewareConfig{Cors: true, JSON: true}))
	// just for test
	Router.GET("/api/user/all", handler.Middleware(user.GetAllUser, handler.MiddlewareConfig{Cors: true, JSON: true}))
	// Router.POST("/api/user/update", handler.Middleware(user.UpdateInfo, handler.MiddlewareConfig{Cors: true, JSON: true}))

	// memo about
	// Router.GET("/api/:id/", memo.GetMemoByID)
	Router.GET("/api/memo/all", handler.Middleware(memo.GetAllMemos, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/memo/new", handler.Middleware(memo.CreateMemo, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/memo/update", handler.Middleware(memo.UpdateMemo, handler.MiddlewareConfig{Cors: true, JSON: true}))
	Router.POST("/api/memo/delete", handler.Middleware(memo.DeleteMemo, handler.MiddlewareConfig{Cors: true, JSON: true}))

	Router.NotFound = http.HandlerFunc(handler.NotFoundHandler)
}
