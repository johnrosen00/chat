package handlers

import (
	"chat/server/gateway/models/users"
	"chat/server/gateway/sessions"
)

//HandlerContext contains pointers to structs necessary for user authorization and state tracking.
type HandlerContext struct {
	Key          string
	SessionStore sessions.Store
	UserStore    users.Store
	SocketStore  *SocketStore
}
