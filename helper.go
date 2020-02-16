package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

type BakeryAndLocale struct {
	Bakery
	Locality
}

func newBakeries(b []Bakery, l []Locality) []BakeryAndLocale {
	var bl []BakeryAndLocale
	for k := range b {
		t := BakeryAndLocale{Bakery: b[k], Locality: l[k]}
		bl = append(bl, t)
	}
	return bl
}

func getSharedBakery(db *gorm.DB, agentID int) []BakeryAndLocale {
	/*
		get bakeryshare from tblbakeryshare
		query bakeries table where fldflouragentno = ?
		// CHECK if this association is correct.
		//FIXME make preload instead of this hacky way
	*/
	var bs BakeryShare
	var baker []Bakery
	// var l []Locality

	// FldFlourAgentNo
	db.Table("tblbakeryshare").Find(&bs, "FldFlourAgentNo = ?", agentID)

	db.Table("tblbakery").Find(&baker, "fldbakeryno = ?", bs.FldBakeryNo)
	ids := func(r []Bakery) []int {
		var i []int
		for _, v := range r {
			i = append(i, v.FldLocalityNo)
		}
		return i
	}(baker)

	// db.Table("tbllocality").Find(&l, "fldlocalityno in (?)", ids)
	var res []BakeryAndLocale
	db.Raw(`
		SELECT
		tb.*, tc.FldCityName, tl.FldLocalityName, ts.FldStateName, tn.FldNeighborhoodName
		FROM TblBakery tb
			INNER JOIN TblCity tc on tc.FldCityNo = tb.FldCityNo
			INNER JOIN TblLocality tl on tl.FldLocalityNo = tb.FldLocalityNo
			INNER JOIN TblState ts on ts.FldStateNo = tb.FldStateNo
			INNER JOIN TblNeighborhood tn on tn.FldNeighborhoodNo = tb.FldNeighborhoodNo
		where tc.FldLocalityNo in (?) 
`, ids).Scan(&res)
	// b := newBakeries(baker, l)
	return res
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

//geo query data according for each Bakery
/* it should work like this:
- it should query city / locality / neighborhood / admin
*/
func geo(db *gorm.DB, agent int, data Geo) []Address {

	var res []Address

	builder := "tblbakery.fldcityno = ?"
	if data.State > 0 {
		builder += " AND tblbakery.fldstateno = ?"
	}
	if data.Neighborhood > 0 {
		builder += "AND tblbakery.FldNeighborhoodNo = ?"
	}
	log.Printf("the data is: %#v", data)

	q := db.Table("tblbakery").Select("tblbakery.*, tc.FldCityName, tl.FldLocalityName, ts.FldStateName, tn.FldNeighborhoodName").Joins(`INNER JOIN TblCity tc on tc.FldCityNo = tblbakery.FldCityNo
	 		INNER JOIN TblLocality tl on tl.FldLocalityNo = tblbakery.FldLocalityNo
			INNER JOIN TblState ts on ts.FldStateNo = tblbakery.FldStateNo
			INNER JOIN TblNeighborhood tn on tn.FldNeighborhoodNo = tblbakery.FldNeighborhoodNo`)

	q.Where(builder, data.City, data.State, data.Neighborhood).Scan(&res)
	log.Printf("the data is: %#v", res)
	return res
}

//Geo this field
type Geo struct {
	City         int
	Locality     int
	Neighborhood int
	Admin        int
	State        int
}

type ListGeo struct {
	City         []NameID
	Neighborhood []NameID
	State        []NameID
	Locality     []NameID
	Admin        []NameID
}

type NameID struct {
	ID   int
	Name string
}

func getID(r *http.Request, param string) int {
	q := r.URL.Query().Get(param)
	qq, _ := strconv.Atoi(q)
	return qq
}

func marshalAddresses(a []Address) []byte {
	d, _ := json.Marshal(&a)
	return d
}
