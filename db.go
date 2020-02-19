package main

import (
	"encoding/json"
	"log"

	"github.com/jinzhu/gorm"
)

//User users table
/*
FldUserNo	FldFullTable	FldUserTable	FldPassword	FldUserType	FLdImage
FldDisabled	FldStateNo	FldLocalityNo	FldCityNo	FldNeighborhoodNo
FldSecurityLevel	FldUpdateDate

// TODO encrypt the password

User
A user can be of different types, Agent, distributor or Grinder
*/
type User struct {
	FldUserNo         int         `gorm:"column:FldUserNo" json:"FldUserNo"`
	FldFullName       string      `gorm:"column:FldFullName" json:"FldFullName"`
	FldFullTable      string      `gorm:"column:FldFullTable" json:"FldFullTable"`
	FldUserTable      string      `gorm:"column:FldUserTable" json:"FldUserTable"`
	FldPassword       string      `gorm:"column:FldPassword" json:"-"`
	FldUserType       int         `gorm:"column:FldUserType" json:"FldUserType"`
	FldImage          interface{} `gorm:"column:FldImage" json:"FldImage"`
	FldDisabled       bool        `gorm:"column:FldDisabled" json:"FldDisabled"`
	FldStateNo        int         `gorm:"column:FldStateNo" json:"FldStateNo"`
	FldLocaliyNo      int         `gorm:"column:FldLocaliyNo" json:"FldLocaliyNo"`
	FldCityNo         int         `gorm:"column:FldCityNo" json:"FldCityNo"`
	FldNeighborhoodNo int         `gorm:"column:FldNeighborhoodNo" json:"FldNeighborhoodNo"`
	FldSecurityLevel  int         `gorm:"column:FldSecurityLevel" json:"FldSecurityLevel"`
	FldUpdateDate     string      `gorm:"column:FldUpdateDate" json:"FldUpdateDate"`
	FldSystemNo       int         `gorm:"column:FldSystemNo" json:"FldSystemNo"`
	FldUserName       string      `gorm:"column:FldUserName" json:"FldUserName"`
}

//getID is supposed to returns user id which will be used throughout the system
// it should map to agentid, bakeryid, grinderid, etc
func (u User) getID() int {
	return u.FldSystemNo
}

func (u User) marshal() []byte {
	data, _ := json.Marshal(&u)
	return data
}
func checkPassword(password string, u User) bool {
	return password == u.FldPassword
}

func getUser(db *gorm.DB, username string) (bool, User) {
	var user User
	if err := db.Table("tblusers").Find(&user, "fldusername = ?", username).Error; err != nil {
		return false, user
	} else {
		return true, user
	}
}

func getProfile(db *gorm.DB, user User) (bool, UserProfile) {

	var res UserProfile
	sys := user.FldSystemNo

	switch user.FldUserType {
	case 6: // case bakery
		db.Table("tblusers").Select("tblusers.*, tb.FldPhone").Joins("INNER JOIN TblBakery tb on tb.FldSystemNo = tblusers.FldSystemNo").Where("tb.fldsystemno = ?", sys).Scan(&res)
	case 7: // case agent
		db.Table("tblusers").Select("tblusers.*, tb.FldPhone").Joins("INNER JOIN tblagent tb on tb.FldSystemNo = tblusers.FldSystemNo")
	case 8: // case ginder
		db.Table("tblusers").Select("tblusers.*, tb.FldPhone").Joins("INNER JOIN tblgrinder tb on tb.FldSystemNo = tblusers.FldSystemNo")
	case 9: // case importer
		db.Table("tblusers").Select("tblusers.*, tb.FldPhone").Joins("INNER JOIN tblimporter tb on tb.FldSystemNo = tblusers.FldSystemNo")

	}
	return true, res
}

//Bakery model
/*
FldBakeryNo	FldBakeryTable	FldIsActive	FldStateNo	FldLocalityNo	FldCityNo
	FldNeighborhoodNo	FldContactTable	FldPhone	FldEmail	FldAddress
		FldVolume	FldLong	FldLat	FldUserNo	FldLogNo	FldUpdateDate
		FldImage	FldNFCBakeryNo
*/
type Bakery struct {
	FldBakeryNo       int    `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldBakeryTable    string `gorm:"column:FldBakeryName" json:"FldBakeryTable,omitempty"`
	FldIsActive       bool   `gorm:"column:FldIsActive" json:"FldIsActive,omitempty"`
	FldStateNo        int    `gorm:"column:FldStateNo" json:"FldStateNo,omitempty"`
	FldLocalityNo     int    `gorm:"column:FldLocalityNo" json:"FldLocalityNo,omitempty"`
	FldCityNo         int    `gorm:"column:FldCityNo" json:"FldCityNo,omitempty"`
	FldNeighborhoodNo int    `gorm:"column:FldNeighborhoodNo" json:"FldNeighborhoodNo,omitempty"`
	FldContactName    string `gorm:"column:FldContactName" json:"FldContactName,omitempty"`
	FldPhone          string `gorm:"column:FldPhone" json:"FldPhone,omitempty"`
	FldEmail          string `gorm:"column:FldEmail" json:"FldEmail,omitempty"`
	FldAddress        string `gorm:"column:FldAddress" json:"FldAddress,omitempty"`
	FldVolume         int    `gorm:"column:FldVolume" json:"FldVolume,omitempty"` // FIXME type
	FldLong           string `gorm:"column:FldLong" json:"FldLong,omitempty"`
	FldLat            string `gorm:"column:FldLat" json:"FldLat,omitempty"`
	FldUserNo         int    `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"` // is this a foreignkey?
	FldLogNo          int    `gorm:"column:FldLogNo" json:"FldLogNo,omitempty"`
	FldUpdateDate     string `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
	FldImage          string `gorm:"column:FldImage" json:"FldImage,omitempty"` // is this bytes, blob, etc?
	FldNFCBakeryNo    int    `gorm:"column:FldNFCBakeryNo" json:"FldNFCBakeryNo,omitempty"`
}

func (b Bakery) TableName() string {
	return "tblbakery"
}

func (b Bakery) getAll(db *gorm.DB, data Geo) []BakeryAndLocale {
	var res []BakeryAndLocale
	db.Raw(`
		SELECT
		tb.*, tc.FldCityName, tl.FldLocalityName, ts.FldStateName, tn.FldNeighborhoodName
		FROM TblBakery tb
			INNER JOIN TblCity tc on tc.FldCityNo = tb.FldCityNo
			INNER JOIN TblLocality tl on tl.FldLocalityNo = tb.FldLocalityNo
			INNER JOIN TblState ts on ts.FldStateNo = tb.FldStateNo
			INNER JOIN TblNeighborhood tn on tn.FldNeighborhoodNo = tb.FldNeighborhoodNo
			where tb.FldStateNo = ? AND tb.FldCityNo = ? AND tb.FldLocalityNo = ?`, data.State, data.City, data.Locality).Scan(&res)
	return res
}

func (b Bakery) getMarshaled(db *gorm.DB, data Geo) []byte {
	d, _ := json.Marshal(b.getAll(db, data))
	return d
}

//FlourAgent
/*
- Flour Agent App
	• Record Received Flour from Grinder [TblFlourAgentReceive] [ Use TblFlourAgentShare as a lookup]
	• Record Distributed Flours to Bakery [TblFlourAgentDistribute] [Use TblBakerShare as a lookup]

FldFlourAgentNo	FldFlourAgentName	FldIsActive	FldStateNo	FldContactName	FldPhone	FldEmail
FldAddress	FldVolume	FldLong	FldLat	FldUserNo	FldLogNo	FldUpdateDate


*/
type FlourAgent struct {
	FldFlourAgentNo   int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldFlourAgentName string  `gorm:"column:FldFlourAgentName" json:"FldFlourAgentName,omitempty"`
	FldIsActive       bool    `gorm:"column:FldIsActive" json:"FldIsActive,omitempty"`
	FldStateNo        int     `gorm:"column:FldStateNo" json:"FldStateNo,omitempty"`
	FldContactName    string  `gorm:"column:FldContactName" json:"FldContactName,omitempty"`
	FldPhone          string  `gorm:"column:FldPhone" json:"FldPhone,omitempty"`
	FldEmail          string  `gorm:"column:FldEmail" json:"FldEmail,omitempty"`
	FldAddress        string  `gorm:"column:FldAddress" json:"FldAddress,omitempty"`
	FldVolume         float32 `gorm:"column:FldVolume" json:"FldVolume,omitempty"`
	FldLong           string  `gorm:"column:FldLong" json:"FldLong,omitempty"`
	FldLat            string  `gorm:"column:FldLat" json:"FldLat,omitempty"`
	FldUserNo         int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
	FldLogNo          int     `gorm:"column:FldLogNo" json:"FldLogNo,omitempty"`
	FldUpdateDate     string  `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

func (a FlourAgent) getAgents(bakeryID int, db *gorm.DB) ([]FlourAgent, error) {
	// BakeryShare
	// FldBakeryNo FldFlourAgentNo FldShareAmount FldFrequency FldIsSelected
	// agent no = 2 for testing

	var f []FlourAgent
	var ids []int
	var s []BakeryShare

	if err := db.Table("tblbakeryshare").Find(&s, "fldbakeryno = ?", bakeryID).Error; err != nil {
		return []FlourAgent{}, err
	}

	log.Printf("the data is: %v", s)
	ids = getInts(s)
	if err := db.Table("tblflouragent").Find(&f, "fldflouragentno in (?)", ids).Error; err != nil {
		return []FlourAgent{}, err
	}
	return f, nil
}

func getInts(b []BakeryShare) []int {
	var a []int
	for _, v := range b {
		a = append(a, v.FldFlourAgentNo)

	}
	return a

}

func (a FlourAgent) marshal() []byte {
	d, _ := json.Marshal(&a)
	return d
}

//FlourAgentReceive
/*
FldFlourAgentReceiveNo	FldDate	FldFlourAgentNo	FldGrinderNo	FldQuantity	FldUnitPrice	FldTotalAmount
FldRefNo	FldNFCFlourAgentReceiveNo	FldNFCStatusNo	FldNFCNote	FldUserNo	FldUpdateDate
*/
type FlourAgentReceive struct {
	FldFlourAgentReceiveNo    int     `gorm:"column:FldFlourAgentReceiveNo" json:"FldFlourAgentReceiveNo,omitempty"`
	FldDate                   string  `gorm:"column:FldDate" json:"FldDate,omitempty"` //FIXME mssql werid smalldatetime bug
	FldFlourAgentNo           int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldGrinderNo              int     `gorm:"column:FldGrinderNo" json:"FldGrinderNo,omitempty"`
	FldQuantity               float32 `gorm:"column:FldQuantity" json:"FldQuantity,omitempty"`
	FldUnitPrice              float32 `gorm:"column:FldUnitPrice" json:"FldUnitPrice,omitempty"`
	FldTotalAmount            float32 `gorm:"column:FldTotalAmount" json:"FldTotalAmount,omitempty"`
	FldRefNo                  string  `gorm:"column:FldRefNo" json:"FldRefNo,omitempty"`
	FldNFCFlourAgentReceiveNo int     `gorm:"column:FldNFCFlourAgentReceiveNo" json:"FldNFCFlourAgentReceiveNo,omitempty"`
	FldNFCStatusNo            int     `gorm:"column:FldNFCStatusNo" json:"FldNFCStatusNo,omitempty"`
	FldNFCNote                string  `gorm:"column:FldNFCNote" json:"FldNFCNote,omitempty"`
	FldUserNo                 int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
	FldUpdateDate             string  `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

func (f FlourAgentReceive) validateReceive() bool {
	if f.FldFlourAgentNo != 0 || f.FldFlourAgentReceiveNo != 0 || f.FldGrinderNo != 0 {
		return true
	}
	// NOT implemented...
	return false
}

func (f FlourAgentReceive) submit(db *gorm.DB) error {
	// template 1900-01-01T00:00:00
	log.Printf("the datetime in flouragent is: %v", f.FldDate)
	// db.Exec("UPDATE orders SET shipped_at=? WHERE id IN (?)", time.Now(), []int64{11,22,33})
	// db.Exec("insert into tblflouragentreceive (fldflouragentno, fldflouragentreceiveno, fldgrinderno, flddate) values (?, ?, ?, ?)", f.FldFlourAgentNo, f.FldFlourAgentReceiveNo, f.FldGrinderNo, f.FldDate)
	if err := db.Table("tblflouragentreceive").Create(&f).Error; err != nil {
		return err
	}
	return nil
}

//getAll gets all data from tblflouragentreceive
func (FlourAgentReceive) getAll(db *gorm.DB) ([]FlourAgentReceive, error) {
	var f []FlourAgentReceive
	if err := db.Table("tblflouragentreceive").Find(&f).Error; err != nil {
		log.Printf("error in flouragentreceive getAll is: %v", err)
		return f, nil
	}
	return f, nil
}

//TableName sets FlourAgentReceive table name to its equivalent sql server name
func (FlourAgentReceive) TableName() string {
	return "tblflouragentreceive"
}

//FlourAgentShare
/*FldFlourAgentNo	FldGrinderNo	FldShareAmount	FldFrequency	FldIsSelected
 */
type FlourAgentShare struct {
	FldFlourAgentNo int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldGrinderNo    int     `gorm:"column:FldGrinderNo" json:"FldGrinderNo,omitempty"`
	FldShareAmount  float32 `gorm:"column:FldShareAmount" json:"FldShareAmount,omitempty"`
	FldFrequency    string  `gorm:"column:FldFrequency" json:"FldFrequency,omitempty"` //FIXME check the field type
	FldIsSelected   bool    `gorm:"column:FldIsSelected" json:"FldIsSelected,omitempty"`
}

//TableName sets FlourAgentShare struct to its equivalent name in the sql server
func (FlourAgentShare) TableName() string {
	return "tblflouragentshare"
}

//FloorAgentDistribute
/*FldFlourAgentDistributeNo	FldDate	FldFlourAgentNo	FldBakeryNo	FldQuantity	FldUnitPrice	FldTotalAmount
FldRefNo	FldNFCFlourBakeryReceiveNo	FldNFCFlourAgentDistributeNo	FldNFCStatusNo	FldNFCNote	FldUserNo	FldUpdateDate
*/
type FlourAgentDistribute struct {
	FldFlourAgentDistributeNo    int     `gorm:"column:FldFlourAgentDistributeNo" json:"FldFlourAgentDistributeNo,omitempty"`
	FldDate                      string  `gorm:"column:FldDate" json:"FldDate,omitempty"`
	FldFlourAgentNo              int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldBakeryNo                  int     `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldQuantity                  float32 `gorm:"column:FldQuantity" json:"FldQuantity,omitempty"`
	FldUnitPrice                 float32 `gorm:"column:FldUnitPrice" json:"FldUnitPrice,omitempty"`
	FldTotalAmount               float32 `gorm:"column:FldTotalAmount" json:"FldTotalAmount,omitempty"`
	FldRefNo                     string  `gorm:"column:FldRefNo" json:"FldRefNo,omitempty"`
	FldNFCFlourBakeryReceiveNo   int     `gorm:"column:FldNFCFlourBakeryReceiveNo" json:"FldNFCFlourBakeryReceiveNo,omitempty"`
	FldNFCFlourAgentDistributeNo int     `gorm:"column:FldNFCFlourAgentDistributeNo" json:"FldNFCFlourAgentDistributeNo,omitempty"`
	FldNFCStatusNo               int     `gorm:"column:FldNFCStatusNo" json:"FldNFCStatusNo,omitempty"`
	FldNFCNote                   string  `gorm:"column:FldNFCNote" json:"FldNFCNote,omitempty"`
	FldUserNo                    int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
}

//TableName sets FlourAgentDistribute struct to its equivalent name in the sql server
func (FlourAgentDistribute) TableName() string {
	return "tblflouragentdistribute"
}

func (f FlourAgentDistribute) validate() bool {
	if f.FldFlourAgentNo != 0 || f.FldFlourAgentDistributeNo != 0 || f.FldBakeryNo != 0 || f.FldQuantity != 0 {
		return true
	}
	return false
}

func (f FlourAgentDistribute) submit(db *gorm.DB) error {
	// template 1900-01-01T00:00:00
	log.Printf("the datetime in flouragent is: %v", f.FldDate)
	// db.Exec("UPDATE orders SET shipped_at=? WHERE id IN (?)", time.Now(), []int64{11,22,33})
	// db.Exec("insert into tblflouragentreceive (fldflouragentno, fldflouragentreceiveno, fldgrinderno, flddate) values (?, ?, ?, ?)", f.FldFlourAgentNo, f.FldFlourAgentReceiveNo, f.FldGrinderNo, f.FldDate)
	if err := db.Table("tblflouragentdistribute").Create(&f).Error; err != nil {
		return err
	}
	return nil
}

//BakeryShare
/*
FldBakeryNo FldFlourAgentNo FldShareAmount FldFrequency FldIsSelected
*/
type BakeryShare struct {
	FldBakeryNo     int     `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldFlourAgentNo int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldShareAmount  float32 `gorm:"column:FldShareAmount" json:"FldShareAmount,omitempty"`
	FldFrequency    string  `gorm:"column:FldFrequency" json:"FldFrequency,omitempty"` //FIXME check the field type
	FldIsSelected   bool    `gorm:"column:FldIsSelected" json:"FldIsSelected,omitempty"`
}

//Grinder
/*FldGrinderNo	FldGrinderName	FldIsActive	FldStateNo	FldContactName	FldPhone	FldEmail
FldAddress	FldVolume	FldUserNo	FldLogNo	FldUpdateDate
*/
type Grinder struct {
	FldGrinderNo   int     `gorm:"column:FldGrinderNo" json:"FldGrinderNo,omitempty"`
	FldGrinderName string  `gorm:"column:FldGrinderName" json:"FldGrinderName,omitempty"`
	FldIsActive    bool    `gorm:"column:FldIsActive" json:"FldIsActive,omitempty"`
	FldStateNo     int     `gorm:"column:FldStateNo" json:"FldStateNo,omitempty"`
	FldContactName string  `gorm:"column:FldContactName" json:"FldContactName,omitempty"`
	FldPhone       string  `gorm:"column:FldPhone" json:"FldPhone,omitempty"`
	FldEmail       string  `gorm:"column:FldEmail" json:"FldEmail,omitempty"`
	FldAddress     string  `gorm:"column:FldAddress" json:"FldAddress,omitempty"`
	FldVolume      float32 `gorm:"column:FldVolume" json:"FldVolume,omitempty"`
	FldUserNo      int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
	FldLogNo       int     `gorm:"column:FldLogNo" json:"FldLogNo,omitempty"` // what is this? HELPNEEDED
	FldUpdateDate  string  `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

func (g Grinder) marshal() []byte {
	d, _ := json.Marshal(&g)
	return d
}

func marshalGrinders(g []Grinder) []byte {
	d, _ := json.Marshal(&g)
	return d
}

func marshalFloursRecv(f []FlourAgentReceive) []byte {
	d, _ := json.Marshal(&f)
	return d
}

func marshalBakeries(b []Bakery) []byte {
	d, _ := json.Marshal(&b)
	return d
}

//marshallBakeriesWithLocale
func marshalBakeriesWithLocale(b []BakeryAndLocale) []byte {
	d, _ := json.Marshal(&b)
	return d
}

// Bakery Tables
/*
FldFlourBakeryReceiveNo	FldDate	FldBakeryNo	FldFlourAgentNo	FldQuantity	FldUnitPrice
FldTotalAmount	FldRefNo	FldDriverName	FldCarPlateNo	FldFlourAgentDistributeNo
FldNFCFlourBakeryReceiveNo	FldNFCStatusNo	FldNFCNote	FldUserNo	FldUpdateDate
*/
type BakeryFlourReceive struct {
	FldFlourBakeryReceiveNo    int     `gorm:"column:FldFlourBakeryReceiveNo" json:"FldFlourBakeryReceiveNo,omitempty"`
	FldFlourAgentDistributeNo  int     `gorm:"column:FldFlourAgentDistributeNo" json:"FldFlourAgentDistributeNo,omitempty"`
	FldDate                    string  `gorm:"column:FldDate" json:"FldDate,omitempty"`
	FldFlourAgentNo            int     `gorm:"column:FldFlourAgentNo" json:"FldFlourAgentNo,omitempty"`
	FldBakeryNo                int     `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldQuantity                float32 `gorm:"column:FldQuantity" json:"FldQuantity,omitempty"`
	FldUnitPrice               float32 `gorm:"column:FldUnitPrice" json:"FldUnitPrice,omitempty"`
	FldTotalAmount             float32 `gorm:"column:FldTotalAmount" json:"FldTotalAmount,omitempty"`
	FldRefNo                   string  `gorm:"column:FldRefNo" json:"FldRefNo,omitempty"`
	FldNFCFlourBakeryReceiveNo int     `gorm:"column:FldNFCFlourBakeryReceiveNo" json:"FldNFCFlourBakeryReceiveNo,omitempty"`
	FldNFCStatusNo             int     `gorm:"column:FldNFCStatusNo" json:"FldNFCStatusNo,omitempty"`
	FldNFCNote                 string  `gorm:"column:FldNFCNote" json:"FldNFCNote,omitempty"`
	FldUserNo                  int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`

	// bakery specific fields
	FldDriverName string `gorm:"column:FldDriverName" json:"FldDriverName,omitempty"`
	FldCarPlateNo string `gorm:"column:FldCarPlateNo" json:"FldCarPlateNo,omitempty"`
	FldUpdateDate string `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

//TableName sets FlourAgentDistribute struct to its equivalent name in the sql server
func (BakeryFlourReceive) TableName() string {
	return "tblflourbakeryreceive"
}

//FIXME what are the fields to validate?
func (f BakeryFlourReceive) validate() bool {
	if f.FldFlourAgentNo != 0 || f.FldFlourAgentDistributeNo != 0 || f.FldBakeryNo != 0 || f.FldQuantity != 0 || f.FldTotalAmount != 0 {
		return true
	}
	return false
}

func (f BakeryFlourReceive) submit(db *gorm.DB) error {
	// Record Received Flour from Flour Agent [TblFlourBakeryReceive] [Use TblBakeyShare as lookup]
	/*
		get bakeryshare from tblbakeryshare
		submit to tableflourbakeryreceive
		FldNFCFlourBakeryReceiveNo
		FldFlourAgentDistributeNo
	*/
	//FIXME it only get's the table now, it doesn't really commit anything yet
	if err := db.Table("TblFlourBakeryReceive").Create(&f).Error; err != nil {
		return err
	}

	return nil
}

func (f BakeryFlourReceive) populate(db *gorm.DB, agentID int) BakeryFlourReceive {
	// Record Received Flour from Flour Agent [TblFlourBakeryReceive] [Use TblBakeryShare as lookup]
	/*
		get bakeryshare from tblbakeryshare
		submit to tableflourbakeryreceive
	*/
	// var b Bakery
	// fixme this function is very ugly
	var bs BakeryShare
	// get backery agent from backery share
	db.Table("tblbakeryshare").Find(&bs, "fldbakeryno = ?", agentID) // fixme check errors
	// fixme what are the fields we want to get?
	f.FldFlourAgentDistributeNo = bs.FldFlourAgentNo
	f.FldFlourAgentNo = bs.FldFlourAgentNo

	return f
}

type FlourBaking struct {
	//FldFlourBakingNo	FldDate	FldBakeryNo	FldWorkingStatusNo	FldQuantity	FldNote	FldLocalityCheck
	// FldLocalityUserNo	FldLocalityNote	FldSecurityCheck	FldSecurityUserNo	FldSecurityNote
	// FldGovernmentalCheck	FldGovermentalUserNo	FldGovernmentalNote	FldCommunityCheck
	// FldComuunityUserNo	FldCommunityNote	FldNFCFlourBakingNo	FldNFCStatusNo	FldNFCNote
	// FldUserNo	FldUpdateDate

	FldFlourBakingNo     int     `gorm:"column:FldFlourBakingNo" json:"FldFlourBakingNo,omitempty"`
	FldDate              string  `gorm:"column:FldDate" json:"FldDate,omitempty"`
	FldBakeryNo          int     `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldWorkingStatusNo   int     `gorm:"column:FldWorkingStatusNo" json:"FldWorkingStatusNo,omitempty"`
	FldQuantity          float32 `gorm:"column:FldQuantity" json:"FldQuantity,omitempty"`
	FldNote              string  `gorm:"column:FldNote" json:"FldNote,omitempty"`
	FldLocalityCheck     float32 `gorm:"column:FldLocalityCheck" json:"FldLocalityCheck,omitempty"`
	FldLocalityUserNo    int     `gorm:"column:FldLocalityUserNo" json:"FldLocalityUserNo,omitempty"`
	FldLocalityNote      string  `gorm:"column:FldLocalityNote" json:"FldLocalityNote,omitempty"`
	FldSecurityCheck     float32 `gorm:"column:FldSecurityCheck" json:"FldSecurityCheck,omitempty"`
	FldSecurityUserNo    int     `gorm:"column:FldSecurityUserNo" json:"FldSecurityUserNo,omitempty"`
	FldSecurityNote      string  `gorm:"column:FldSecurityNote" json:"FldSecurityNote,omitempty"`
	FldGovernmentalCheck float32 `gorm:"column:FldGovernmentalCheck" json:"FldGovernmentalCheck,omitempty"`
	FldGovermentalUserNo int     `gorm:"column:FldGovermentalUserNo" json:"FldGovermentalUserNo,omitempty"`
	FldGovernmentalNote  string  `gorm:"column:FldGovernmentalNote" json:"FldGovernmentalNote,omitempty"`
	FldCommunityCheck    float32 `gorm:"column:FldCommunityCheck" json:"FldCommunityCheck,omitempty"`
	FldComuunityUserNo   int     `gorm:"column:FldComuunityUserNo" json:"FldComuunityUserNo,omitempty"`
	FldCommunityNote     string  `gorm:"column:FldCommunityNote" json:"FldCommunityNote,omitempty"`
	FldNFCFlourBakingNo  int     `gorm:"column:FldNFCFlourBakingNo" json:"FldNFCFlourBakingNo,omitempty"`
	FldNFCStatusNo       int     `gorm:"column:FldNFCStatusNo" json:"FldNFCStatusNo,omitempty"`
	FldNFCNote           string  `gorm:"column:FldNFCNote" json:"FldNFCNote,omitempty"`
	FldUserNo            int     `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
	FldUpdateDate        string  `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

func (f FlourBaking) validate() bool {
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]
	if f.FldDate == "" || f.FldBakeryNo == 0 || f.FldQuantity == 0 {
		return false
	}
	return true
}

// validateAuditor validate auditor submissions data
func (f FlourBaking) validateAuditor() bool {
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]
	if !f.validate() {
		return false
	}

	// perform extra validation here as needed.
	return true
}

func (f FlourBaking) populate(agentID int) FlourBaking {
	f.FldBakeryNo = agentID
	return f
}

func (f FlourBaking) submit(db *gorm.DB) error {
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]
	if err := db.Table("tblflourbaking").Create(&f).Error; err != nil {
		return err
	}
	return nil
}

// TblBakeryAudit
/*FldBakeryAuditNo	FldDate	FldBakeyNo	FldAuditBy	FldAuditType
FldAuditStatusNo	FldNote	FldAuditResponseNo	FldNFCBakeryAuditNo
FldNFCStatusNo	FldNFCNote	FldUserNo	FldUpdateDate */

//BakeryAudit for reporting issues on bakeries
type BakeryAudit struct {
	FldBakeryAuditNo    int    `gorm:"column:FldBakeryAuditNo" json:"FldBakeryAuditNo,omitempty"`
	FldDate             string `gorm:"column:FldDate" json:"FldDate,omitempty"`
	FldBakeyNo          int    `gorm:"column:FldBakeyNo" json:"FldBakeyNo,omitempty"`
	FldAuditBy          int    `gorm:"column:FldAuditBy" json:"FldAuditBy,omitempty"`
	FldAuditType        int    `gorm:"column:FldAuditType" json:"FldAuditType,omitempty"`
	FldAuditStatusNo    int    `gorm:"column:FldAuditStatusNo" json:"FldAuditStatusNo,omitempty"`
	FldNote             string `gorm:"column:FldNote" json:"FldNote,omitempty"`
	FldAuditResponseNo  int    `gorm:"column:FldAuditResponseNo" json:"FldAuditResponseNo,omitempty"`
	FldNFCBakeryAuditNo int    `gorm:"column:FldNFCBakeryAuditNo" json:"FldNFCBakeryAuditNo,omitempty"`
	FldNFCStatusNo      int    `gorm:"column:FldNFCStatusNo" json:"FldNFCStatusNo,omitempty"`
	FldNFCNote          string `gorm:"column:FldNFCNote" json:"FldNFCNote,omitempty"`
	FldUserNo           int    `gorm:"column:FldUserNo" json:"FldUserNo,omitempty"`
	FldUpdateDate       string `gorm:"column:FldUpdateDate" json:"FldUpdateDate,omitempty"`
}

//Name sets table name to match what is in the DB
func (BakeryAudit) Name() string {
	return "tblbakeryaudit"
}

func (b BakeryAudit) submit(db *gorm.DB) error {
	if err := db.Table("tblbakeryaudit").Create(&b).Error; err != nil {
		return err
	}
	return nil
}

func (b BakeryAudit) validate() bool {
	// TODO make validations here
	return true
}

func (b BakeryAudit) populate(agentID int) BakeryAudit {
	// TODO make validations here
	b.FldUserNo = agentID
	return b
}

func (b BakeryAudit) getBakeries(db *gorm.DB, agent int) []BakeryAudit {
	var res []BakeryAudit

	db.Table("tblbakeryaudit").Where("fldbakeryauditno = ?", agent).Find(&res)
	return res
}

//AuditStatus table for inquiring complains
type AuditStatus struct {
	FldAuditStatusNo   int    `gorm:"primary_key;column:FldAuditStatusNo;" json:"FldAuditStatusNo,omitempty"`
	FldAuditStatusName string `gorm:"column:FldAuditStatusName" json:"FldAuditStatusName,omitempty"`
}

func (a AuditStatus) migrate(db *gorm.DB) {
	db.AutoMigrate(&a)
}

//TableName overrides default gorm naming
func (a AuditStatus) TableName() string {
	return "tblauditstatus"
}

func (a AuditStatus) marshal() []byte {
	d, _ := json.Marshal(&a)
	return d
}

func (a AuditStatus) getMarshalled(db *gorm.DB) ([]byte, error) {
	if err := db.Table("tblauditstatus").Find(&a).Error; err != nil {
		return nil, err
	}
	d, _ := json.Marshal(&a)
	return d, nil
}

func getAllComplains(db *gorm.DB) []AuditStatus {
	var r []AuditStatus
	db.Table("tblauditstatus").Find(&r)
	return r
}

func (a AuditStatus) generate(db *gorm.DB) {
	db.AutoMigrate(&a)
	db.DropTable(&a)
	db.CreateTable(&a)
	data := []string{"not_available", "empty_bakery", "missing_ele"}
	for _, v := range data {
		a.FldAuditStatusName = v
		db.Create(&a)
	}
}

//Locality provides locality info
type Locality struct {
	FldLocalityNo       int    `gorm:"primary_key;column:FldLocalityNo" json:"-"`
	FldLocalityName     string `gorm:"column:FldLocalityName" json:"FldLocalityName,omitempty"`
	FldCityName         string `gorm:"column:FldCityName" json:"FldCityName,omitempty"`
	FldStateName        string `gorm:"column:FldStateName" json:"FldStateName,omitempty"`
	FldNeighborhoodName string `gorm:"column:FldNeighborhoodName" json:"FldNeighborhoodName,omitempty"`
}

func (l Locality) TableName() string {
	return "tbllocality"
}

type Address struct {
	FldLocalityName string `gorm:"column:FldLocalityName" json:"FldLocalityName,omitempty"`
	FldLocalityNo   int    `gorm:"column:FldLocalityNo" json:"FldLocalityNo,omitempty"`

	FldCityName string `gorm:"column:FldCityName" json:"FldCityName,omitempty"`
	FldCityNo   int    `gorm:"column:FldCityNo" json:"FldCityNo,omitempty"`

	FldStateName string `gorm:"column:FldStateName" json:"FldStateName,omitempty"`
	FldStateNo   int    `gorm:"column:FldStateNo" json:"FldStateNo,omitempty"`

	FldNeighborhoodName string `gorm:"column:FldNeighborhoodName" json:"FldNeighborhoodName,omitempty"`
	FldNeighborhoodNo   int    `gorm:"column:FldNeighborhoodNo" json:"FldNeighborhoodNo,omitempty"`

	FldAdminName string `gorm:"column:FldAdminName" json:"FldAdminName,omitempty"`
	FldAdminNo   int    `gorm:"column:FldAdminNo" json:"FldAdminNo,omitempty"`
}

func (a Address) marshal() []byte {
	d, _ := json.Marshal(&a)
	return d
}

type UserProfile struct {
	User
	FldPhone string `gorm:"column:FldPhone" json:"FldPhone"`
}

func (u UserProfile) marshal() []byte {
	d, _ := json.Marshal(&u)
	return d
}
