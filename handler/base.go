package handler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MiddlewareConfig config
type MiddlewareConfig struct {
	Cors bool
	JSON bool
}

// Middleware the middleware function
func Middleware(fn httprouter.Handle, config MiddlewareConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("here", config)
		if config.Cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if config.JSON {
			w.Header().Set("Content-Type", "application/json")
		}

		fn(w, r, p)
	}
}
