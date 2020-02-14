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
	FldUserNo         int         `gorm:"column:FldUserNo"`
	FldFullName       string      `gorm:"column:FldFullName"`
	FldFullTable      string      `gorm:"column:FldFullTable"`
	FldUserTable      string      `gorm:"column:FldUserTable"`
	FldPassword       string      `gorm:"column:FldPassword" json:"-"`
	FldUserType       int         `gorm:"column:FldUserType"`
	FldImage          interface{} `gorm:"column:FldImage"`
	FldDisabled       bool        `gorm:"column:FldDisabled"`
	FldStateNo        int         `gorm:"column:FldStateNo"`
	FldLocaliyNo      int         `gorm:"column:FldLocaliyNo"`
	FldCityNo         int         `gorm:"column:FldCityNo"`
	FldNeighborhoodNo int         `gorm:"column:FldNeighborhoodNo"`
	FldSecurityLevel  int         `gorm:"column:FldSecurityLevel"`
	FldUpdateDate     string      `gorm:"column:FldUpdateDate"`
	FldSystemNo       int         `gorm:"column:FldSystemNo"`
	FldUserName       string      `gorm:"column:FldUserName"`
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

func getUser(db *gorm.DB, userTable string) (bool, User) {
	var user User
	if err := db.Table("tblusers").Find(&user, "flduserTable = ?", userTable).Error; err != nil {
		return false, user
	} else {
		return true, user
	}
}

//Bakery model
/*
FldBakeryNo	FldBakeryTable	FldIsActive	FldStateNo	FldLocalityNo	FldCityNo
	FldNeighborhoodNo	FldContactTable	FldPhone	FldEmail	FldAddress
		FldVolume	FldLong	FldLat	FldUserNo	FldLogNo	FldUpdateDate
		FldImage	FldNFCBakeryNo
*/
type Bakery struct {
	FldBakeryNo       int    `gorm:"column:FldBakeryNo"`
	FldBakeryTable    string `gorm:"column:FldBakeryName"`
	FldIsActive       bool   `gorm:"column:FldIsActive"`
	FldStateNo        int    `gorm:"column:FldStateNo"`
	FldLocalityNo     int    `gorm:"column:FldLocalityNo"`
	FldCityNo         int    `gorm:"column:FldCityNo"`
	FldNeighborhoodNo int    `gorm:"column:FldNeighborhoodNo"`
	FldContactName    string `gorm:"column:FldContactName"`
	FldPhone          string `gorm:"column:FldPhone"`
	FldEmail          string `gorm:"column:FldEmail"`
	FldAddress        string `gorm:"column:FldAddress"`
	FldVolume         int    `gorm:"column:FldVolume"` // FIXME type
	FldLong           string `gorm:"column:FldLong"`
	FldLat            string `gorm:"column:FldLat"`
	FldUserNo         int    `gorm:"column:FldUserNo"` // is this a foreignkey?
	FldLogNo          int    `gorm:"column:FldLogNo"`
	FldUpdateDate     string `gorm:"column:FldUpdateDate"`
	FldImage          string `gorm:"column:FldImage"` // is this bytes, blob, etc?
	FldNFCBakeryNo    int    `gorm:"column:FldNFCBakeryNo"`
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
	FldFlourAgentNo   int     `gorm:"column:FldFlourAgentNo"`
	FldFlourAgentName string  `gorm:"column:FldFlourAgentName"`
	FldIsActive       bool    `gorm:"column:FldIsActive"`
	FldStateNo        int     `gorm:"column:FldStateNo"`
	FldContactName    string  `gorm:"column:FldContactName"`
	FldPhone          string  `gorm:"column:FldPhone"`
	FldEmail          string  `gorm:"column:FldEmail"`
	FldAddress        string  `gorm:"column:FldAddress"`
	FldVolume         float32 `gorm:"column:FldVolume"`
	FldLong           string  `gorm:"column:FldLong"`
	FldLat            string  `gorm:"column:FldLat"`
	FldUserNo         int     `gorm:"column:FldUserNo"`
	FldLogNo          int     `gorm:"column:FldLogNo"`
	FldUpdateDate     string  `gorm:"column:FldUpdateDate"`
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
	FldFlourAgentReceiveNo    int     `gorm:"column:FldFlourAgentReceiveNo"`
	FldDate                   string  `gorm:"column:FldDate"` //FIXME mssql werid smalldatetime bug
	FldFlourAgentNo           int     `gorm:"column:FldFlourAgentNo"`
	FldGrinderNo              int     `gorm:"column:FldGrinderNo"`
	FldQuantity               float32 `gorm:"column:FldQuantity"`
	FldUnitPrice              float32 `gorm:"column:FldUnitPrice"`
	FldTotalAmount            float32 `gorm:"column:FldTotalAmount"`
	FldRefNo                  int     `gorm:"column:FldRefNo"`
	FldNFCFlourAgentReceiveNo int     `gorm:"column:FldNFCFlourAgentReceiveNo"`
	FldNFCStatusNo            int     `gorm:"column:FldNFCStatusNo"`
	FldNFCNote                string  `gorm:"column:FldNFCNote"`
	FldUserNo                 int     `gorm:"column:FldUserNo"`
	FldUpdateDate             string  `gorm:"column:FldUpdateDate"`
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
	FldFlourAgentNo int     `gorm:"column:FldFlourAgentNo"`
	FldGrinderNo    int     `gorm:"column:FldGrinderNo"`
	FldShareAmount  float32 `gorm:"column:FldShareAmount"`
	FldFrequency    string  `gorm:"column:FldFrequency"` //FIXME check the field type
	FldIsSelected   bool    `gorm:"column:FldIsSelected"`
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
	FldFlourAgentDistributeNo    int     `gorm:"column:FldFlourAgentDistributeNo"`
	FldDate                      string  `gorm:"column:FldDate"`
	FldFlourAgentNo              int     `gorm:"column:FldFlourAgentNo"`
	FldBakeryNo                  int     `gorm:"column:FldBakeryNo"`
	FldQuantity                  float32 `gorm:"column:FldQuantity"`
	FldUnitPrice                 float32 `gorm:"column:FldUnitPrice"`
	FldTotalAmount               float32 `gorm:"column:FldTotalAmount"`
	FldRefNo                     int     `gorm:"column:FldRefNo"`
	FldNFCFlourBakeryReceiveNo   int     `gorm:"column:FldNFCFlourBakeryReceiveNo"`
	FldNFCFlourAgentDistributeNo int     `gorm:"column:FldNFCFlourAgentDistributeNo"`
	FldNFCStatusNo               int     `gorm:"column:FldNFCStatusNo"`
	FldNFCNote                   string  `gorm:"column:FldNFCNote"`
	FldUserNo                    int     `gorm:"column:FldUserNo"`
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
	FldBakeryNo     int     `gorm:"column:FldBakeryNo"`
	FldFlourAgentNo int     `gorm:"column:FldFlourAgentNo"`
	FldShareAmount  float32 `gorm:"column:FldShareAmount"`
	FldFrequency    string  `gorm:"column:FldFrequency"` //FIXME check the field type
	FldIsSelected   bool    `gorm:"column:FldIsSelected"`
}

//Grinder
/*FldGrinderNo	FldGrinderName	FldIsActive	FldStateNo	FldContactName	FldPhone	FldEmail
FldAddress	FldVolume	FldUserNo	FldLogNo	FldUpdateDate
*/
type Grinder struct {
	FldGrinderNo   int     `gorm:"column:FldGrinderNo"`
	FldGrinderName string  `gorm:"column:FldGrinderName"`
	FldIsActive    bool    `gorm:"column:FldIsActive"`
	FldStateNo     int     `gorm:"column:FldStateNo"`
	FldContactName string  `gorm:"column:FldContactName"`
	FldPhone       string  `gorm:"column:FldPhone"`
	FldEmail       string  `gorm:"column:FldEmail"`
	FldAddress     string  `gorm:"column:FldAddress"`
	FldVolume      float32 `gorm:"column:FldVolume"`
	FldUserNo      int     `gorm:"column:FldUserNo"`
	FldLogNo       int     `gorm:"column:FldLogNo"` // what is this? HELPNEEDED
	FldUpdateDate  string  `gorm:"column:FldUpdateDate"`
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
	FldFlourBakeryReceiveNo    int     `gorm:"column:FldFlourBakeryReceiveNo"`
	FldFlourAgentDistributeNo  int     `gorm:"column:FldFlourAgentDistributeNo"`
	FldDate                    string  `gorm:"column:FldDate"`
	FldFlourAgentNo            int     `gorm:"column:FldFlourAgentNo"`
	FldBakeryNo                int     `gorm:"column:FldBakeryNo"`
	FldQuantity                float32 `gorm:"column:FldQuantity"`
	FldUnitPrice               float32 `gorm:"column:FldUnitPrice"`
	FldTotalAmount             float32 `gorm:"column:FldTotalAmount"`
	FldRefNo                   int     `gorm:"column:FldRefNo"`
	FldNFCFlourBakeryReceiveNo int     `gorm:"column:FldNFCFlourBakeryReceiveNo"`
	FldNFCStatusNo             int     `gorm:"column:FldNFCStatusNo"`
	FldNFCNote                 string  `gorm:"column:FldNFCNote"`
	FldUserNo                  int     `gorm:"column:FldUserNo"`

	// bakery specific fields
	FldDriverName string `gorm:"column:FldDriverName"`
	FldCarPlateNo string `gorm:"column:FldCarPlateNo"`
	FldUpdateDate string `gorm:"column:FldUpdateDate"`
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

	FldFlourBakingNo     int     `gorm:"column:FldFlourBakingNo"`
	FldDate              string  `gorm:"column:FldDate"`
	FldBakeryNo          int     `gorm:"column:FldBakeryNo"`
	FldWorkingStatusNo   int     `gorm:"column:FldWorkingStatusNo"`
	FldQuantity          float32 `gorm:"column:FldQuantity"`
	FldNote              string  `gorm:"column:FldNote"`
	FldLocalityCheck     float32 `gorm:"column:FldLocalityCheck"`
	FldLocalityUserNo    int     `gorm:"column:FldLocalityUserNo"`
	FldLocalityNote      string  `gorm:"column:FldLocalityNote"`
	FldSecurityCheck     float32 `gorm:"column:FldSecurityCheck"`
	FldSecurityUserNo    int     `gorm:"column:FldSecurityUserNo"`
	FldSecurityNote      string  `gorm:"column:FldSecurityNote"`
	FldGovernmentalCheck float32 `gorm:"column:FldGovernmentalCheck"`
	FldGovermentalUserNo int     `gorm:"column:FldGovermentalUserNo"`
	FldGovernmentalNote  int     `gorm:"column:FldGovernmentalNote"`
	FldCommunityCheck    float32 `gorm:"column:FldCommunityCheck"`
	FldComuunityUserNo   int     `gorm:"column:FldComuunityUserNo"`
	FldCommunityNote     int     `gorm:"column:FldCommunityNote"`
	FldNFCFlourBakingNo  int     `gorm:"column:FldNFCFlourBakingNo"`
	FldNFCStatusNo       int     `gorm:"column:FldNFCStatusNo"`
	FldNFCNote           string  `gorm:"column:FldNFCNote"`
	FldUserNo            int     `gorm:"column:FldUserNo"`
	FldUpdateDate        string  `gorm:"column:FldUpdateDate"`
}

func (f FlourBaking) validate() bool {
	// Record Baked Flour [TblFlourBaking]  [Set FldDate,FldBakeryNo, FldQunatity, FldNote]
	if f.FldDate == "" || f.FldBakeryNo == 0 || f.FldQuantity == 0 || f.FldNote == "" {
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
	FldBakeryAuditNo    int    `gorm:"column:FldBakeryAuditNo"`
	FldDate             string `gorm:"column:FldDate"`
	FldBakeyNo          int    `gorm:"column:FldBakeyNo"`
	FldAuditBy          int    `gorm:"column:FldAuditBy"`
	FldAuditType        int    `gorm:"column:FldAuditType"`
	FldAuditStatusNo    int    `gorm:"column:FldAuditStatusNo"`
	FldNote             string `gorm:"column:FldNote"`
	FldAuditResponseNo  int    `gorm:"column:FldAuditResponseNo"`
	FldNFCBakeryAuditNo int    `gorm:"column:FldNFCBakeryAuditNo"`
	FldNFCStatusNo      int    `gorm:"column:FldNFCStatusNo"`
	FldNFCNote          string `gorm:"column:FldNFCNote"`
	FldUserNo           int    `gorm:"column:FldUserNo"`
	FldUpdateDate       string `gorm:"column:FldUpdateDate"`
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

//AuditStatus table for inquiring complains
type AuditStatus struct {
	FldAuditStatusNo   int    `gorm:"primary_key;column:FldAuditStatusNo;"`
	FldAuditStatusName string `gorm:"column:FldAuditStatusName"`
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

// Table Locality

type Locality struct {
	FldLocalityNo   int    `gorm:"primary_key;column:FldLocalityNo"`
	FldLocalityName string `gorm:"column:FldLocalityName"`
}

func (l Locality) TableName() string {
	return "tbllocality"
}
