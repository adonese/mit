package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

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

func getGrinderFromAgent(db *gorm.DB, agentID int) (bool, Grinder) {
	var grinder Grinder
	var s FlourAgentShare
	err := db.Table("tblflouragentshare").Find(&s, "fldflouragentno = ?", agentID).Error
	if err != nil {
		return false, Grinder{}
	}
	log.Printf("grinder no is: %v", s.FldGrinderNo)
	err = db.Table("tblgrinder").Find(&grinder, "fldgrinderno = ?", s.FldGrinderNo).Error
	if err != nil {
		return false, Grinder{}
	}
	return true, grinder
}

func getBakeryFromAgent(db *gorm.DB, agentID int) Grinder {
	var grinder Grinder
	var s FlourAgentShare
	db.Table("tblflooragentshare").Find(&s, "fldagentno = ?", agentID)
	db.Table("tblgrinder").Find(&grinder, "fldgrinderno = ?", s.FldGrinderNo)
	return grinder
}
