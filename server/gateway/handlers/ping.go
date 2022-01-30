package handlers

import (
	"chat/server/models/users"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
)

//PingDB pings databases for error handling.
func (cx *HandlerContext) PingDB(w http.ResponseWriter, r *http.Request) {

	userStore := &users.MySQLStore{}
	err := errors.New("")
	msg := "no error"
	if userStore.DB, err = sql.Open("mysql", os.Getenv("DSN")); err != nil {
		msg = fmt.Sprintf("error opening database: %v\n", err)
		w.Write([]byte(msg))
	}

	defer userStore.DB.Close()

	if err = userStore.DB.Ping(); err != nil {
		msg = fmt.Sprintf("error pinging database: %v\n", err)
		w.Write([]byte(msg))
	}

	w.Write([]byte("UserStore all good"))

}
