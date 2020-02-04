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

	http.HandleFunc("/login", login)
	http.ListenAndServe(":8091", nil)
}
