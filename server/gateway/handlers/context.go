package handlers

import (
	"chat/server/gateway/sessions"
	"chat/server/models/db"
)

//HandlerContext contains pointers to structs necessary for user authorization and state tracking.
type HandlerContext struct {
	Key          string
	SessionStore sessions.Store
	Data         *db.Connection
	//SocketStore  *SocketStore
}
