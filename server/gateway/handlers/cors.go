package handlers

import "net/http"

//WrappedCORSHandler is a middleware handler that wraps a handler to manipulate
type WrappedCORSHandler struct {
	handler http.Handler
}

//NewWrappedCORSHandler constructs a new WrappedCORSHandler middleware handler
func NewWrappedCORSHandler(handlerToWrap http.Handler) *WrappedCORSHandler {
	return &WrappedCORSHandler{handlerToWrap}
}

//ServeHTTP handles the request by adding the response header
func (h *WrappedCORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add the header
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.Header().Add("Access-Control-Max-Age", "600")
	//call the wrapped handler

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	h.handler.ServeHTTP(w, r)
}
