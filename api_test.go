package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_login(t *testing.T) {
	// r := login()
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
