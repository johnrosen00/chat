package handlers

// import (
// 	"chat/server/gateway/sessions"
// 	"chat/server/models/rabbitpayload"
// 	//"chat/server/models/rabbitpayload"
// 	"encoding/json"
// 	"net/http"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// const (
// 	// TextMessage denotes a text data message. The text message payload is
// 	// interpreted as UTF-8 encoded text data.
// 	TextMessage = 1

// 	// BinaryMessage denotes a binary data message.
// 	BinaryMessage = 2

// 	// CloseMessage denotes a close control message. The optional message
// 	// payload contains a numeric code and text. Use the FormatCloseMessage
// 	// function to format a close message payload.
// 	CloseMessage = 8

// 	// PingMessage denotes a ping control message. The optional message payload
// 	// is UTF-8 encoded text.
// 	PingMessage = 9

// 	// PongMessage denotes a pong control message. The optional message payload
// 	// is UTF-8 encoded text.
// 	PongMessage = 10
// )

// //TODO: add a handler that upgrades clients to a WebSocket connection
// //and adds that to a list of WebSockets to notify when events are
// //read from the RabbitMQ server. Remember to synchronize changes
// //to this list, as handlers are called concurrently from multiple
// //goroutines.
// //and adds that to a list of WebSockets to notify when events are
// //read from the RabbitMQ server. Remember to synchronize changes
// //to this list, as handlers are called concurrently from multiple
// //goroutines.

// //WSUpgradeHandler upgrades client connection to ws://
// func (cx *HandlerContext) WSUpgradeHandler(w http.ResponseWriter, r *http.Request) {

// 	//check whether current user is authenticated
// 	currentSess := &SessionState{}
// 	if _, err := sessions.GetState(r, cx.Key, cx.SessionStore, currentSess); err != nil {
// 		http.Error(w, "Please sign in.", http.StatusUnauthorized)
// 		return
// 	}
// 	id := currentSess.User.ID

// 	origin := r.Header.Get("Origin")
// 	if origin != "https://johnrosen.me" && origin != "https://api.johnrosen.me" {
// 		http.Error(w, "Websocket Connection Refused", 403)
// 	}

// 	//upgrade client connection
// 	upgrader := websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			// if request origin header bad {
// 			// 	return false
// 			// }
// 			return true
// 		},
// 	}

// 	//add current connection to websocket data structure.
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		http.Error(w, "Error upgrading to websocket.", http.StatusBadRequest)
// 		return
// 	}
// 	cx.SocketStore.HandleConnection(conn, id)

// 	go func() {
// 		cx.connectionRoutine(id)
// 	}()
// }

// //HandleConnection inserts a connection, and creates a goroutine to send and receive messages.
// func (s *SocketStore) HandleConnection(conn *websocket.Conn, id int64) error {
// 	s.Lock.Lock()
// 	// insert socket connection
// 	s.Connections[id] = conn
// 	s.Lock.Unlock()
// 	//starts
// 	return nil
// }

// //RemoveConnection removes a websocket connection
// func (s *SocketStore) removeConnection(id int64) error {
// 	s.Lock.Lock()

// 	//stop socket reading goroutine
// 	s.Connections[id].Close()
// 	// delete socket connection
// 	delete(s.Connections, id)

// 	s.Lock.Unlock()

// 	return nil
// }

// //Goroutine to run for each active websocket connection
// //listens to when to close connection
// func (cx *HandlerContext) connectionRoutine(id int64) {
// 	conn := cx.SocketStore.Connections[id]
// 	defer conn.Close()
// 	defer cx.SocketStore.removeConnection(id)
// 	for {
// 		messageType, _, err := conn.ReadMessage()

// 		if messageType == CloseMessage {
// 			cx.SocketStore.removeConnection(id)
// 			break
// 		}

// 		if err != nil {
// 			cx.SocketStore.removeConnection(id)
// 			break
// 		}
// 	}
// 	// cleanup
// }

// //socketstore struct, defined here instead of models because go was throwing a fit

// //SocketStore is a store for socket connections, with a built in Mutex lock
// type SocketStore struct {
// 	Connections map[int64]*websocket.Conn
// 	Lock        sync.Mutex
// }

// //InitSocketStore initializes a socketstore
// func InitSocketStore() *SocketStore {
// 	this := &SocketStore{}
// 	this.Connections = make(map[int64]*websocket.Conn)

// 	return this
// }

// //Methods to write to users.

// //WriteToUser writes a message to a single user
// func (s *SocketStore) WriteToUser(message []byte, id int64) error {
// 	s.Connections[id].WriteJSON(message)
// 	return nil
// }

// //WriteToUsers writes a message to list of users
// func (s *SocketStore) WriteToUsers(message []byte, ids []int64) error {
// 	for i := int64(0); i < int64(len(ids)); i++ {
// 		if s.Connections[i] != nil {
// 			if err := s.WriteToUser(message, ids[i]); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// //WriteToAllUsers writes a message to all users
// func (s *SocketStore) WriteToAllUsers(message []byte) error {
// 	for i := int64(0); i < int64(len(s.Connections)); i++ {
// 		if err := s.WriteToUser(message, i); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// //GatewayConsumeHandler handles consumption in the gateway
// func (s *SocketStore) GatewayConsumeHandler(body []byte) error {
// 	payload := &rabbitpayload.Payload{}
// 	if err := json.Unmarshal(body, payload); err != nil {
// 		return err
// 	}

// 	if payload.UserIDs == nil {
// 		s.WriteToAllUsers(payload.Body)
// 	} else {
// 		s.WriteToUsers(payload.Body, payload.UserIDs)
// 	}
// 	return nil
// }
