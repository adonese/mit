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

	// log.Printf("the user profile is: %#v", d)

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

	_, d := getProfile(db, u)

	w.WriteHeader(http.StatusOK)
	w.Write(d.marshal())
}

func generateToken() {
	//TODO
}

func logout(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func refreshToken(w http.ResponseWriter, r *http.Request) {}

//getGrinderHandler gets associated grinders to specific agent, using agent ID (provided in url query params)
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
	ok, g := getGrinderFromAgent1(db, agentID)
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

//getBakeries get associated bakeries to this agent
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
	w.Write(marshalBakeriesWithLocale(b))
}

//getBakeries get associated bakeries to this agent
//experiemtnal api to get agent with bakeries
func agentBakeries(w http.ResponseWriter, r *http.Request) {

	id := getID(r, "agent")
	c := getID(r, "city")
	l := getID(r, "locality")
	n := getID(r, "neighborhood")
	a := getID(r, "admin")
	s := getID(r, "state")

	db := getEngine()
	data := Geo{Locality: l, City: c, Admin: a, Neighborhood: n, State: s}

	w.Header().Add("content-type", "application/json")
	b := getAgentSharedBakeries(db, id, data)
	w.WriteHeader(http.StatusOK)
	w.Write(marshalBakeriesWithLocale(b))
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

func bakeryAgentsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]

	// get bakeryid
	id := r.URL.Query().Get("agent")
	bakeryID, _ := strconv.Atoi(id)

	w.Header().Add("content-type", "application/json")
	db := getEngine()
	defer db.Close()

	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Malformed request", Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var ag FlourAgent
	a, err := ag.getAgents(bakeryID, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ve := validationError{Message: err.Error(), Code: "agents_flour_retrieval_err"}
		w.Write(ve.marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshalFlourAgents(a))

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

//auditorCheckHandler reports the flour quantity by an auditor
// this needs further thinking
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

		ve := validationError{Message: err.Error(), Code: "empty_request_body"}
		w.Write(ve.marshal())
		return
	}
	defer r.Body.Close()

	var f flourData
	if err = json.Unmarshal(req, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ve := validationError{Message: err.Error(), Code: "marshalling_error"}
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

	var b FlourBaking
	modBakery := b.populateAuditors(bakeryID, f)
	if err := modBakery.submit(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ve := validationError{Message: err.Error(), Code: "server_error"}
		w.Write(ve.marshal())
		return
	}

	s := success{Result: "ok"}

	w.WriteHeader(http.StatusOK)

	w.Write(s.marshal())
	return
}

//auditorGetBaked let the auditor get the baked amount from a bakery
func auditorGetBaked(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]

	w.Header().Add("content-type", "application/json")
	db := getEngine()
	state := getID(r, "state")
	admin := getID(r, "admin")
	locality := getID(r, "locality")

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	f := flourData{Start: start, End: end, State: state, Admin: admin, Locality: locality}

	log.Printf("the flour data is: %#v", f)
	if ok := f.validate(); !ok {
		w.WriteHeader(http.StatusBadRequest)

		ve := validationError{Message: "Some fields are missing", Code: "missing_fields"}
		w.Write(ve.marshal())
		return
	}

	// populate bakery data

	// FIXME this part is extremely ugly

	var b FlourBaking
	geo := Geo{State: f.State, Locality: f.Locality, Admin: f.Admin}
	d := b.getBakedMarshaled(db, geo, f.Start, f.End)
	w.Write(d)
	w.WriteHeader(http.StatusOK)
	return
}

//violationHandler reports an incident in a respective bakery
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

//getComplains a drop down list to get complains from
func getComplains(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db := getEngine()
	var a AuditStatus

	a.migrate(db)
	complains := getAllComplains(db)
	res := marshalAuditStatus(complains)

	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func auditorBakeries(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("agent")
	agent, _ := strconv.Atoi(id)
	state := getID(r, "state")
	locality := getID(r, "locality")
	admin := getID(r, "admin")

	geo := Geo{State: state, Locality: locality, Admin: admin}

	w.Header().Add("content-type", "application/json")
	db := getEngine()
	defer db.Close()

	rr := BakeryAudit{}.filterBakeries(db, agent, geo)
	w.Write(marshalBakeries(rr))
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

//generateComplains: FOR TESTING purposes only. Not exposed via any api
func generateComplains(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	db := getEngine()
	var a AuditStatus

	a.migrate(db)
	a.generate(db)
	w.WriteHeader(http.StatusOK)
}

//getBakeries get all bakeries
func getAllBakeries(w http.ResponseWriter, r *http.Request) {
	// return all bakeries. do the filtering later
	w.Header().Add("content-type", "application/json")

	// id := getID(r, "agent")
	c := getID(r, "city")
	l := getID(r, "locality")
	n := getID(r, "neighborhood")
	a := getID(r, "admin")
	s := getID(r, "state")
	data := Geo{Locality: l, City: c, Admin: a, Neighborhood: n, State: s}

	db := getEngine()
	d := Bakery{}.getMarshaled(db, data)
	w.Write(d)
}

func getLocations(w http.ResponseWriter, r *http.Request) {
	// return all bakeries. do the filtering later
	w.Header().Add("content-type", "application/json")

	id := getID(r, "agent")
	c := getID(r, "city")
	l := getID(r, "locality")
	n := getID(r, "neighborhood")
	a := getID(r, "admin")
	s := getID(r, "state")

	db := getEngine()
	data := Geo{Locality: l, City: c, Admin: a, Neighborhood: n, State: s}
	d := geo(db, id, data)
	log.Printf("the data to be printed is: %v", d)
	w.Write(marshalAddresses(d))
}

//getLocalities gets
func getLocalities(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	id := getID(r, "agent")
	c := getID(r, "city")
	l := getID(r, "locality")
	n := getID(r, "neighborhood")
	a := getID(r, "admin")
	s := getID(r, "state")

	db := getEngine()
	data := Geo{Locality: l, City: c, Admin: a, Neighborhood: n, State: s}
	d := getCustomLocations(db, id, data)
	log.Printf("the data to be printed is: %v", d)
	w.Write(marshalAddresses(d))
}
