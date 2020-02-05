package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func submitFlourHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var f FlourAgentReceive
	json.Unmarshal(req, &f)

	if ok := f.validateReceive(); !ok {
		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ve.marshal())
	}
	f.submit(db)
	w.WriteHeader(http.StatusOK)
	s := success{Result: "ok"}
	w.Write(s.marshal())
}

//getSubmittedFlourHandler
/* This is highly advised to be only admin's view */
func getSubmittedFlourHandler(w http.ResponseWriter, r *http.Request) {
	// todo
	w.Header().Add("content-type", "application/json")
	db := getEngine()
	var flour FlourAgentReceive
	f, err := flour.getAll(db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: err.Error(), Code: "records_not_found"}
		w.Write(ve.marshal())
		return
	}

	res := marshalFloursRecv(f)
	w.Write(res)
	w.WriteHeader(http.StatusOK)
	return

}

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

func generateToken() {}

func logout(w http.ResponseWriter, r *http.Request) {}

func refreshToken(w http.ResponseWriter, r *http.Request) {}

func getGrinderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db := getEngine()
	id := r.URL.Query().Get("agent")
	agentID, _ := strconv.Atoi(id)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Agent is not provided", Code: "missing_agent"}
		w.Write(ve.marshal())
		return
	}
	ok, g := getGrinderFromAgent(db, agentID)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Grinder is not available", Code: "missing_grinder"}
		w.Write(ve.marshal())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshalGrinders(g))
}
