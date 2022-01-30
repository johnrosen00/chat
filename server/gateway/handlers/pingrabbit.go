package handlers

import (
	"net/http"

	"github.com/streadway/amqp"
)

//PingRabbit checks to see if rabbit is ok
func (cx *HandlerContext) PingRabbit(w http.ResponseWriter, r *http.Request) {
	_, err := amqp.Dial("amqp://guest:guest@queue:5672/")

	if err != nil {
		w.Write([]byte("no heartbeat"))
	} else {
		w.Write([]byte("connection ok"))
	}
}
