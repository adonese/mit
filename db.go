package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

//User users table
/*
FldUserNo	FldFullName	FldUserName	FldPassword	FldUserType	FLdImage
FldDisabled	FldStateNo	FldLocaliyNo	FldCityNo	FldNeighborhoodNo
FldSecurityLevel	FldUpdateDate

// TODO encrypt the password

User
A user can be of different types, Agent, distributor or Grinder
*/
type User struct {
	FldUserNo         int         `gorm:"column:FldUserNo"`
	FldFullName       string      `gorm:"column:FldFullName"`
	FldUserName       string      `gorm:"column:FldUserName"`
	FldPassword       string      `gorm:"column:FldPassword" json:"-"`
	FldUserType       int         `gorm:"column:FldUserType"`
	FldImage          interface{} `gorm:"column:FldImage"`
	FldDisabled       bool        `gorm:"column:FldDisabled"`
	FldStateNo        int         `gorm:"column:FldStateNo"`
	FldLocaliyNo      int         `gorm:"column:FldLocaliyNo"`
	FldCityNo         int         `gorm:"column:FldCityNo"`
	FldNeighborhoodNo int         `gorm:"column:FldNeighborhoodNo"`
	FldSecurityLevel  int         `gorm:"column:FldSecurityLevel"`
	FldUpdateDate     time.Time   `gorm:"column:FldUpdateDate"`
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

//Bakery model
/*
FldBakeryNo	FldBakeryName	FldIsActive	FldStateNo	FldLocalityNo	FldCityNo
	FldNeighborhoodNo	FldContactName	FldPhone	FldEmail	FldAddress
		FldVolume	FldLong	FldLat	FldUserNo	FldLogNo	FldUpdateDate
		FldImage	FldNFCBakeryNo
*/
type Bakery struct {
	FldBakeryNo       int     `gorm:"column:FldBakeryNo"`
	FldBakeryName     string  `gorm:"column:FldBakeryName"`
	FldIsActive       bool    `gorm:"column:FldIsActive"`
	FldStateNo        int     `gorm:"column:FldStateNo"`
	FldLocalityNo     int     `gorm:"column:FldLocalityNo"`
	FldCityNo         int     `gorm:"column:FldCityNo"`
	FldNeighborhoodNo int     `gorm:"column:FldNeighborhoodNo"`
	FldContactName    string  `gorm:"column:FldContactName"`
	FldPhone          string  `gorm:"column:FldPhone"`
	FldEmail          string  `gorm:"column:FldEmail"`
	FldAddress        string  `gorm:"column:FldAddress"`
	FldVolume         int     `gorm:"column:FldVolume"` // FIXME type
	FldLong           float64 `gorm:"column:FldLong"`
	FldLat            float64 `gorm:"column:FldLat"`
	FldUserNo         int     `gorm:"column:FldUserNo"` // is this a foreignkey?
	FldLogNo          int     `gorm:"column:FldLogNo"`
	FldUpdateDate     string  `gorm:"column:FldUpdateDate"`
	FldImage          string  `gorm:"column:FldImage"` // is this bytes, blob, etc?
	FldNFCBakeryNo    int     `gorm:"column:FldNFCBakeryNo"`
}

//FlourAgent
/*
- Flour Agent App
	• Record Received Flour from Grinder [TblFlourAgentReceive] [ Use TblFlourAgentShare as a lookup]
	• Record Distributed Flours to Bakery [TblFlourAgentDistribute] [Use TblBakerShare as a lookup]

FldFlourAgentNo	FldFlourAgentName	FldIsActive	FldStateNo	FldContactName	FldPhone	FldEmail
FldAddress	FldVolume	FldLong	FldLat	FldUserNo	FldLogNo	FldUpdateDate
*/
type FlourAgent struct{}

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
		return f, nil
	}
	return f, errors.New("data not available")
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
	FldFlourAgentDistributeNo    int       `gorm:"column:FldFlourAgentDistributeNo"`
	FldDate                      time.Time `gorm:"column:FldDate"`
	FldFlourAgentNo              int       `gorm:"column:FldFlourAgentNo"`
	FldBakeryNo                  int       `gorm:"column:FldBakeryNo"`
	FldQuantity                  float32   `gorm:"column:FldQuantity"`
	FldUnitPrice                 float32   `gorm:"column:FldUnitPrice"`
	FldTotalAmount               float32   `gorm:"column:FldTotalAmount"`
	FldRefNo                     int       `gorm:"column:FldRefNo"`
	FldNFCFlourBakeryReceiveNo   int       `gorm:"column:FldNFCFlourBakeryReceiveNo"`
	FldNFCFlourAgentDistributeNo int       `gorm:"column:FldNFCFlourAgentDistributeNo"`
	FldNFCStatusNo               int       `gorm:"column:FldNFCStatusNo"`
	FldNFCNote                   string    `gorm:"column:FldNFCNote"`
	FldUserNo                    int       `gorm:"column:FldUserNo"`
}

//TableName sets FlourAgentDistribute struct to its equivalent name in the sql server
func (FlourAgentDistribute) TableName() string {
	return "tblflouragentdistribute"
}

//Grinder
/*FldGrinderNo	FldGrinderName	FldIsActive	FldStateNo	FldContactName	FldPhone	FldEmail
FldAddress	FldVolume	FldUserNo	FldLogNo	FldUpdateDate
*/
type Grinder struct {
	FldGrinderNo   int       `gorm:"column:FldGrinderNo"`
	FldGrinderName string    `gorm:"column:FldGrinderName"`
	FldIsActive    bool      `gorm:"column:FldIsActive"`
	FldStateNo     int       `gomr:"column:FldStateNo"`
	FldContactName string    `gorm:"column:FldContactName"`
	FldPhone       string    `gorm:"column:FldPhone"`
	FldEmail       string    `gorm:"column:FldEmail"`
	FldAddress     string    `gorm:"column:FldAddress"`
	FldVolume      float32   `gorm:"column:FldVolume"`
	FldUserNo      int       `gorm:"column:FldUserNo"`
	FldLogNo       int       `gorm:"column:FldLogNo"` // what is this? HELPNEEDED
	FldUpdateDate  time.Time `gorm:"column:FldUpdateDate"`
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
