package helper

import (
	"neosmemo/backend/util"
	"net/http"
	"time"
)

// Session session type
type Session struct {
	UserID    string
	SessionID string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// SessionManager SessionID-Session.UserID
//
// Session Manager Design Purpose:
// 1. 用加密的 session text 取代 user_id
// 2. 多端登录（待用 redis/sql 实现）
// 3. and others
//
// NOTE: 暂时存在内存里
var SessionManager = map[string]Session{}

func init() {
	// do nth
}

// GetUserIDFromSession GetUserIDFromSession
func GetUserIDFromSession(r *http.Request) (string, bool) {
	sessionID, err := util.GetKeyValueFromCookie("session_id", r)
	if err != nil {
		return "", false
	}

	session, ok := SessionManager[sessionID]
	if !ok {
		return "", false
	}

	return session.UserID, true
}
