package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	var login Login
	w.Header().Add("content-type", "application/json")
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.Unmarshal(b, &login) // check errors here

	db := getEngine()
	ok, u := getUser(db, login.Username)
	if !ok {
		ve := validationError{Message: "User not found", Code: "user_not_found"}
		data, _ := json.Marshal(&ve)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	} else {
		if ok := checkPassword(login.Password, u); !ok {
			ve := validationError{Message: "Wrong password", Code: "wrong_password"}
			data, _ := json.Marshal(&ve)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u.marshal())
}

func logout(w http.ResponseWriter, r *http.Request) {}

func refreshToken(w http.ResponseWriter, r *http.Request) {}
