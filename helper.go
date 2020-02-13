package main

import (
	"encoding/json"
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

//getGrinderFromAgent it receives a agentID and then fetches it to get
// *its associated flouragent. It uses the flour agent to load its profile
//HELPNEEDED: is these assumptions are correct?
func getGrinderFromAgent(db *gorm.DB, agentID int) (bool, []Grinder) {
	var grinder []Grinder
	var s FlourAgentShare
	err := db.Table("tblflouragentshare").Find(&s, "fldflouragentno = ?", agentID).Error
	if err != nil {
		return false, []Grinder{}
	}

	log.Printf("grinder no is: %v", s.FldGrinderNo)
	err = db.Table("tblgrinder").Where("fldgrinderno = ?", s.FldGrinderNo).Find(&grinder).Error

	log.Printf("grinder list is: %v", grinder)
	if err != nil {
		return false, []Grinder{}
	}
	return true, grinder
}

func getGrinderFromAgent1(db *gorm.DB, agentID int) (bool, []Grinder) {
	var grinder []Grinder
	var s []FlourAgentShare
	err := db.Table("tblflouragentshare").Find(&s, "fldflouragentno = ?", agentID).Error
	if err != nil {
		return false, []Grinder{}
	}
	a := func(a []FlourAgentShare) []int {
		var r []int
		for _, v := range a {
			r = append(r, v.FldGrinderNo)
		}

		return r
	}(s)

	// log.Printf("grinder no is: %v", s.FldGrinderNo)
	err = db.Table("tblgrinder").Where("fldgrinderno in (?)", a).Find(&grinder).Error

	log.Printf("grinder list is: %v", grinder)
	if err != nil {
		return false, []Grinder{}
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

func getSharedBakery(db *gorm.DB, agentID int) []Bakery {
	/*
		get bakeryshare from tblbakeryshare
		query bakeries table where fldflouragentno = ?
		// CHECK if this association is correct.
	*/
	var bs BakeryShare
	var baker []Bakery

	// FldFlourAgentNo
	db.Table("tblbakeryshare").Find(&bs, "FldFlourAgentNo = ?", agentID)
	db.Table("tblbakery").Find(&baker, "fldbakeryno = ?", bs.FldBakeryNo)

	return baker
}

//FIXME
func getAgentFromBakery(db *gorm.DB, bakeryID int) int {
	/*
		get bakeryshare from tblbakeryshare
		submit to tableflourbakeryreceive
	*/
	var bs BakeryShare

	db.Table("tblbakeryshare").Find(&bs, "FldBakeryNo = ?", bakeryID)
	// db.Table("tblbaker").Find(&agent, "fldflouragentno = ?", bs.FldFlourAgentNo)

	return bs.FldFlourAgentNo
}

func marshalFlourAgents(a []FlourAgent) []byte {
	d, _ := json.Marshal(&a)
	return d
}

func marshalAuditStatus(a []AuditStatus) []byte {
	d, _ := json.Marshal(&a)
	return d
}
