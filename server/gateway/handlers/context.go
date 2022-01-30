package handlers

import (
	"chat/server/gateway/sessions"
	"chat/server/models/users"
)

//HandlerContext contains pointers to structs necessary for user authorization and state tracking.
type HandlerContext struct {
	Key          string
	SessionStore sessions.Store
	UserStore    users.Store
	//SocketStore  *SocketStore
}
