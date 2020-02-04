package main

import "github.com/jinzhu/gorm"

// MIT users table
/*
FldUserNo	FldFullName	FldUserName	FldPassword	FldUserType	FLdImage
FldDisabled	FldStateNo	FldLocaliyNo	FldCityNo	FldNeighborhoodNo
FldSecurityLevel	FldUpdateDate

// TODO encrypt the password
*/
type User struct {
	FldUserNo         int    `gorm:"column:FldUserNo"`
	FldFullName       string `gorm:"column:FldFullName"`
	FldPassword       string `gorm:"column:FldPassword"`
	FldUserType       int    `gorm:"column:FldUserType"`
	FldImage          string `gorm:"column:FldImage"`
	FldDisabled       bool   `gorm:"column:FldDisabled"`
	FldStateNo        int    `gorm:"column:FldStateNo"`
	FldLocaliyNo      int    `gorm:"column:FldLocaliyNo"`
	FldCityNo         int    `gorm:"column:FldCityNo"`
	FldNeighborhoodNo int    `gorm:"column:FldNeighborhoodNo"`
	FldSecurityLevel  int    `gorm:"column:FldSecurityLevel"`
	FldUpdateDate     string `gorm:"column:FldUpdateDate"`
}

func (u User) checkPassword(password string) bool {
	if password == u.FldPassword {
		return true
	}
	return false
}

func getUser(db *gorm.DB, username string) (bool, error) {
	db.Table("tblusers").Find(&user, "where ...interface{}")
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
	FldAddress        string  `gorm:"column:FldAddress`
	FldVolume         int     `gorm:"column:FldVolume"` // FIXME type
	FldLong           float64 `gorm:"column:FldLong"`
	FldLat            float64 `gorm:"column:FldLat"`
	FldUserNo         int     `gorm:"column:FldUserNo"` // is this a foreignkey?
	FldLogNo          int     `gorm:"column:FldLogNo"`
	FldUpdateDate     string  `gorm:"column:FldUpdateDate"`
	FldImage          string  `gorm:"column:FldImage"` // is this bytes, blob, etc?
	FldNFCBakeryNo    int     `gorm:"column:FldNFCBakeryNo"`
}
