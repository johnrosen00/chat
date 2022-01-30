package handlers

import (
	"chat/server/gateway/sessions"
	"chat/server/models/users"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//UserHandler POSTS new users to database.
func (cx *HandlerContext) UserHandler(w http.ResponseWriter, r *http.Request) {

	//if not POST, will handle GET in future iterations
	if r.Method != "POST" {
		http.Error(w, "Forbidden method.", http.StatusMethodNotAllowed)
		return
	}

	//if request method is POST

	//check content type;
	//	if content type is not application/json
	//		http.Error(w, "Request body must contain JSON.", http.StatusUnsupportedMediaType)

	if !isJSONctype(r) {
		http.Error(w, "Request body must contain JSON.", http.StatusUnsupportedMediaType)
		return
	}

	//check if format is in format of *users.NewUser
	//json.Unmarshal(jsontounmarshal, NU)
	newUser := &users.NewUser{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(newUser); err != nil {
		http.Error(w, "Incorrect JSON format.", http.StatusBadRequest)
		return
	}

	user, err1 := newUser.ToUser()

	if err1 != nil {
		http.Error(w, "Invalid JSON fields.", http.StatusUnsupportedMediaType)
		return
	}

	user, err1 = cx.UserStore.Insert(user)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	//begin new session

	newSessionState := &SessionState{}
	newSessionState.User = user
	newSessionState.StartTime = time.Now()
	sessions.BeginSession(cx.Key, cx.SessionStore, newSessionState, w)

	if errT := cx.UserStore.Track(r, user.ID, time.Now()); errT != nil {
		http.Error(w, errT.Error(), http.StatusBadRequest)
		return
	}

	//	respond to client
	//		set header type to application/json
	//		set status code to http.StatusCreated (ie 201)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//write response body
	//	encoded as json object
	//	id field contains database assigned PK value
	//	passhash and email automatically omitted
	enc := json.NewEncoder(w)
	if err3 := enc.Encode(user); err3 != nil {
		w.Write([]byte("error encoding to json"))
	}
}

//returns true if current request header = application/json
func isJSONctype(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

//SpecificUserHandler handles requests for a specific user.
//path format: /v1/users/{UserID} ie v1/users/1234
func (cx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "GET" && method != "PATCH" {
		http.Error(w, "Unsupported HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	userPath := r.URL.Path
	userPath = path.Base(userPath)
	u := &users.User{}
	if userPath == "me" {

		currentSessionState := &SessionState{}
		if _, errGetSession := sessions.GetState(r, cx.Key, cx.SessionStore, currentSessionState); errGetSession != nil {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}

		u = currentSessionState.User

	} else {
		id0, err := strconv.Atoi(userPath)
		if err != nil {
			http.Error(w, "Bad path.", http.StatusBadRequest)
			return
		}

		//u = user
		u, err = cx.UserStore.GetByID(int64(id0))

		if err != nil {
			http.Error(w, "ID not found.", http.StatusNotFound)
			return
		}
	}

	id := u.ID

	if r.Method == "PATCH" {

		if !isJSONctype(r) {
			http.Error(w, "Unsupported content-type.", http.StatusUnsupportedMediaType)
		}

		newUpdate := &users.Updates{}
		d := json.NewDecoder(r.Body)

		if err := d.Decode(newUpdate); err != nil {
			http.Error(w, "Incorrect JSON format.", http.StatusUnsupportedMediaType)
			return
		}

		u2, err1 := cx.UserStore.Update(id, newUpdate)

		if err1 != nil {
			http.Error(w, "Invalid update fields.", http.StatusUnsupportedMediaType)
			return
		}
		u = u2
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)

	if err3 := enc.Encode(u); err3 != nil {
		http.Error(w, "Unable to encode to JSON", 404)
	}

}

//SessionsHandler creates new sesssions
func (cx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Forbidden method.", http.StatusMethodNotAllowed)
		return
	}

	if !isJSONctype(r) {
		http.Error(w, "Request body must contain JSON.", http.StatusUnsupportedMediaType)
		return
	}

	//	request body must contain json that can be decoded into a users.Credentials struct
	//	if:
	//	- cannot find user using credentials OR
	//	- invalid json format OR
	//	- can find user using credentials but cannot authenticate.
	//		http.Error(w, "invalid credentials", http.StatusUnauthorized)

	newCredentials := &users.Credentials{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(newCredentials); err != nil {
		bcrypt.CompareHashAndPassword([]byte("g"), []byte(newCredentials.Password))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	u, errAuth := cx.UserStore.GetByEmail(newCredentials.Email)
	if errAuth != nil {
		bcrypt.CompareHashAndPassword([]byte("g"), []byte(newCredentials.Password))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword(u.PassHash, []byte(newCredentials.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	//	start new session.
	newSessionState := &SessionState{}
	newSessionState.User = u
	newSessionState.StartTime = time.Now()
	sessions.BeginSession(cx.Key, cx.SessionStore, newSessionState, w)

	if errT := cx.UserStore.Track(r, u.ID, time.Now()); errT != nil {
		http.Error(w, "error tracking user session.", http.StatusBadRequest)
		return
	}
	//	response body:
	//		content-type header = application/json
	// 		status code http.StatusCreated
	//		json encoded copy of user's profile
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)

	if err3 := enc.Encode(u); err3 != nil {
		http.Error(w, "Unable to encode to JSON", 404)
	}

}

//SpecificSessionsHandler ends a specific user's session. only path permitted: /v1/sessions/mine
func (cx *HandlerContext) SpecificSessionsHandler(w http.ResponseWriter, r *http.Request) {
	//if not DELETE, return
	//	http.Error(w, "Unsupported HTTP method.", http.StatusMethodNotAllowed)
	if r.Method != "DELETE" {
		http.Error(w, "Unsupported HTTP method.", http.StatusMethodNotAllowed)
		return
	}
	//if request method is DELETE
	//	if last path segment != "mine"
	//		http.Error(w, "Access forbidden.", http.StatusForbidden)
	pathSeg := r.URL.Path
	pathSeg = path.Base(pathSeg)

	if pathSeg != "mine" {
		http.Error(w, "Access forbidden.", http.StatusForbidden)
		return
	}

	//	end current session.
	_, err := sessions.EndSession(r, cx.Key, cx.SessionStore)

	if err != nil {
		http.Error(w, "invalid session token", http.StatusUnauthorized)
		return
	}
	//	response body:
	//		plain text message: "signed out"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("signed out"))
}
