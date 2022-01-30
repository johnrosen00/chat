package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//TestCORSHeaderAdd tests CORS middleware
func TestCORSHeaderAdd(t *testing.T) {
	sampleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("error encoding to json"))
	})

	handlerToTest := NewWrappedCORSHandler(sampleHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)
	headers := recorder.Header()
	expected := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
	}
	for i := 0; i < len(expected); i++ {
		if len(headers.Get(expected[i])) == 0 {
			t.Error("FAILED, no header " + expected[i])
		} else {
			t.Error("PASS, header " + expected[i])
		}
	}
}
