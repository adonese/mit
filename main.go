package main

import (
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func main() {
	db := getEngine()
	defer db.Close()
	// if err := db.DB().Ping(); err != nil {
	// 	log.Printf("there is an error: %v", err)
	// }
	// users := getUsersTable(db, "tblusers")
	// log.Printf("the users table are: %v", users)

	http.HandleFunc("/listing", listing)
	http.HandleFunc("/login", login)

	// agent apis
	http.HandleFunc("/get_grinders", getGrinderHandler) // to be compatible with #ISSUE 1
	http.HandleFunc("/get_grinder", getGrinderHandler)
	http.HandleFunc("/submit_flour", submitFlourHandler)
	http.HandleFunc("/_get_flour", getSubmittedFlourHandler)
	http.HandleFunc("/get_bakery", getBakeries)
	http.HandleFunc("/submit_bakery", setDistributedFlours)

	// names are really getting too bad now

	// bakery endpoints
	//TODO get associated agents to bakery
	http.HandleFunc("/bakery_submit", bakerySubmitFlourHandler)
	http.HandleFunc("/record_baked", recordBakedHandler)
	http.HandleFunc("/bakery/agents", bakeryAgentsHandler)
	http.HandleFunc("/bakery/submit", bakerySubmitFlourHandler)
	http.HandleFunc("/bakery/received", recordBakedHandler)
	http.HandleFunc("/bakery/record_received", recordBakedHandler)
	http.HandleFunc("/bakery/baked", recordBakedHandler)

	// auditor handlers
	http.HandleFunc("/auditor/check", auditorCheckHandler)
	http.HandleFunc("/auditor/report", violationHandler)

	http.ListenAndServe(":8091", nil)
}
