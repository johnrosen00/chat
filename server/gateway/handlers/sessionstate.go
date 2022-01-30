package handlers

import (
	"chat/server/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

//SessionState describes a session state. Depends on Users package.
type SessionState struct {
	StartTime time.Time   `json:"starttime"`
	User      *users.User `json:"user"`
}
