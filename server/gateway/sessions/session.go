package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	// create a new SessionID
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("invalidSigningKey")
	}
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	//- save the sessionState to the store
	err = store.Save(sessionID, sessionState)

	if err != nil {
		return InvalidSessionID, err
	}

	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)
	bearer := "Bearer " + sessionID.String()
	w.Header().Add("Authorization", bearer)

	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("invalidSigningKey")
	}
	resp1 := r.Header.Get("Authorization")
	resp2 := r.URL.Query().Get("auth")
	resp := resp1

	if len(resp1) == 0 {
		if len(resp2) == 0 {
			return InvalidSessionID, errors.New("http request does not contain ID")
		}

		resp = resp2
	}

	if strings.Contains(resp, "Bearer") {
		resp = strings.Split(resp, " ")[1]
		resp = strings.TrimSpace(resp)
		if len(resp) == 0 {
			return InvalidSessionID, errors.New("no SID")
		}
	} else {
		return InvalidSessionID, errors.New("Invalid scheme")
	}

	sid, err := ValidateID(resp, signingKey)

	if err != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("invalidSigningKey")
	}
	sid, err := GetSessionID(r, signingKey)

	if err != nil {
		return InvalidSessionID, err
	}

	err = store.Get(sid, sessionState)

	if err != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("invalidSigningKey")
	}
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	err = store.Delete(sessionID)

	if err != nil {
		return InvalidSessionID, err
	}

	return sessionID, nil

}
