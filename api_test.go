package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_login(t *testing.T) {
	// r := login()

	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	wrongData := Login{Username: "mohamed", Password: "my wrong password"}
	correctData := Login{Username: "admin", Password: "admin"}
	tests := []struct {
		name string
		req  Login
		want int
	}{
		{"400 request", Login{}, 400}, {"200 request", correctData, 200}, {"400 request wrong payload", wrongData, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := marshal(tt.req)
			dd := bytes.NewBuffer(d)
			res, err := http.Post(ts.URL, "application/json", dd)
			if err != nil {
				log.Fatal(err)
			}
			_, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getUser() got = %v, want %v", res.StatusCode, tt.want)
			}
		})
	}
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func marshal(l Login) []byte {
	data, _ := json.Marshal(&l)
	return data
}
