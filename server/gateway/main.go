package main

import (
	"chat/server/gateway/handlers"
	"chat/server/gateway/sessions"
	"chat/server/models/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
func main() {

	//summary api stuff
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	//auth api stuff
	sessionKey := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")

	sessionStore := &sessions.RedisStore{}

	sessionStore.SessionDuration, _ = time.ParseDuration("999m")

	sessionStore.Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//DB Connections
	dsn := os.Getenv("DSN")
	conn * db.Connection
	if database, err := sql.Open("mysql", dsn); err != nil {
		fmt.Printf("error opening database: %v\n,", err)
	} else {
		defer database.Close()
		conn := db.InitConnection(database)
	}

	cx := &handlers.HandlerContext{
		Key:          sessionKey,
		SessionStore: sessionStore,
		Data:         conn,
	}

	mux2 := http.NewServeMux()

	mux2.HandleFunc("/v1/sessions", cx.SessionsHandler)
	wrappedMux := handlers.NewWrappedCORSHandler(mux2)

	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
	//for testing purposes
	//log.Fatal(http.ListenAndServe(addr, wrappedMux))
}
