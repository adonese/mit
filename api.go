package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	var login *Login
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.Unmarshal(b, login) // check errors here

}
