package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func submitFlourHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f FlourAgentReceive
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validateReceive(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}
	if err := f.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.Write(s.marshal())
	w.WriteHeader(http.StatusOK)

	return
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

	log.Printf("the login from Zooba is: %v", b)
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

func generateToken() {
	//TODO
}

func logout(w http.ResponseWriter, r *http.Request) {
	//TODO
}

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

func setDistributedFlours(w http.ResponseWriter, r *http.Request) {
	// i have only an agent ID. use table agentbakeryshare

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f FlourAgentDistribute
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validate(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}
	if err := f.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.Write(s.marshal())
	w.WriteHeader(http.StatusOK)

	return
}

func getBakeries(w http.ResponseWriter, r *http.Request) {
	db := getEngine()
	// TODO we need to add more queries over here to geo locate and
	// make the results less
	// since an agent can have tons of places
	// also check table of locales
	agentID := r.URL.Query().Get("agent")
	id, _ := strconv.Atoi(agentID)
	b := getSharedBakery(db, id)
	w.WriteHeader(http.StatusOK)
	w.Write(marshalBakeries(b))
}

//TblFlourBakeryReceive
func bakerySubmitFlourHandler(w http.ResponseWriter, r *http.Request) {
	// i have only an agent ID. use table agentbakeryshare

	// get bakeryid
	id := r.URL.Query().Get("agent")
	bakeryID, _ := strconv.Atoi(id)

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f BakeryFlourReceive
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validate(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}

	// populate bakery data
	modBakery := f.populate(db, bakeryID)
	if err := modBakery.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	s := success{Result: "ok"}
	w.Write(s.marshal())

	return
}

//recordBakedHandler endpoint for baker to record the amount of baked bread
func recordBakedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]

	// get bakeryid
	id := r.URL.Query().Get("agent")
	bakeryID, _ := strconv.Atoi(id)

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f FlourBaking
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validate(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}

	// populate bakery data

	// FIXME this part is extremely ugly

	modBakery := f.populate(bakeryID)
	if err := modBakery.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.Write(s.marshal())
	w.WriteHeader(http.StatusOK)

	return
}

func auditorCheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]

	// get bakeryid
	id := r.URL.Query().Get("agent")
	bakeryID, _ := strconv.Atoi(id)

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f FlourBaking
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validateAuditor(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}

	// populate bakery data

	// FIXME this part is extremely ugly

	modBakery := f.populate(bakeryID)
	if err := modBakery.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.Write(s.marshal())
	w.WriteHeader(http.StatusOK)

	return
}

//TblBakeryAudit
func violationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]

	// get bakeryid
	id := r.URL.Query().Get("agent")
	bakeryID, _ := strconv.Atoi(id)

	w.Header().Add("content-type", "application/json")
	db := getEngine()

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f BakeryAudit
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		log.Printf("the error is: %v", err)
		w.Write(ve.marshal())
		return
	}

	if ok := f.validate(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}

	// populate bakery data

	// FIXME this part is extremely ugly

	modBakery := f.populate(bakeryID)
	if err := modBakery.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.Write(s.marshal())
	w.WriteHeader(http.StatusOK)

	return
}

func listing(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var l Listing
	d := l.GetAll()
	w.Write(d.marshal())
	w.WriteHeader(http.StatusOK)
	return
}
