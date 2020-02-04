package main

import "github.com/jinzhu/gorm"

func getUsersTable(db *gorm.DB, tablename string) []User {
	var output []User
	db.Table(tablename).Find(&output)
	return output

}

// getEngine get instance of db connection pool
//FIXME add error handling and refactor "connection string"
func getEngine() *gorm.DB {
	db, _ := gorm.Open("mssql", "sqlserver://SA:Adonese=1994@197.251.5.78:1433?database=MIT")

	return db
}
