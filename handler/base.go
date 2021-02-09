package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MiddlewareConfig config
type MiddlewareConfig struct {
	Cors bool
	JSON bool
}

// Response res type
type Response struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Succeed       bool        `json:"succeed"`
	Data          interface{} `json:"data"`
}

// Middleware the middleware function.
// fn callback function,
// config middleware config.
func Middleware(fn httprouter.Handle, config MiddlewareConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// NOTE: error handler
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("%v", r)
				ErrorHandler(w, string(msg))
			}
		}()

		if config.Cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if config.JSON {
			w.Header().Set("Content-Type", "application/json")
		}

		fn(w, r, p)
	}
}

// ErrorHandler error handler
// NOTE: 可以再抽象一下
func ErrorHandler(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")

	errorRes := Response{
		StatusCode:    http.StatusBadRequest,
		StatusMessage: message,
		Succeed:       false,
	}

	json.NewEncoder(w).Encode(&errorRes)
}
