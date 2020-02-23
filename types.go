package main

import "encoding/json"

// Login is request body for logging-in
type Login struct {
	Username string
	Password string
}

type validationError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (v validationError) marshal() []byte {
	d, _ := json.Marshal(&v)
	return d
}

type success struct {
	Result string `json:"result"`
}

func (s success) marshal() []byte {
	d, _ := json.Marshal(&s)
	return d
}

type Listing struct {
	Agent   map[string][]string
	Bakery  map[string][]string
	Auditor map[string][]string
}

func (l Listing) GetAll() Listing {

	a := make(map[string][]string)
	a["/get_grinders"] = []string{"GET", "Get all grinders associated with this agent"}
	a["/submit_flour"] = []string{"POST", "Records what the agent has received"}
	a["/get_bakery"] = []string{"GET", "Get associated bakeries to this agent"}
	a["/submit_bakery"] = []string{"POST", "Agent records their submitted flour to bakery. You will need to supply bakery ID from the drop down menu"}

	b := make(map[string][]string)
	b["/bakery/submit"] = []string{"POST", "Bakery submits their received flour", "/bakery_submit"}
	b["/bakery/baked"] = []string{"POST", "Bakery submits their *baked* bread", "/bakery/received", "/bakery/record_received"}

	ad := make(map[string][]string)
	ad["/auditor/check"] = []string{"POST", "Auditor {Security, Community, etc}, write the flour received at the X Bakery"}
	ad["/auditor/baked"] = []string{"POST", "Auditor {Security, Community, etc}, write the flour received at the X Bakery"}
	l.Agent = a
	l.Bakery = b
	l.Auditor = ad
	return l
}

func (l Listing) marshal() []byte {
	b, _ := json.Marshal(&l)
	return b
}

type bakeryAndAuditor struct {
	FldBakeryNo   int    `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldBakeryName string `gorm:"column:FldBakeryName" json:"FldBakeryName,omitempty"`
}

type flourData struct {
	// FldLocalityCheck, FldLoclityUserno, FldLocalitynote
	FldLocalityCheck float32 `json:"FldLocalityCheck,omitempty"`
	FldLoclityUserno int     `json:"FldLoclityUserno,omitempty"`
	FldLocalitynote  string  `json:"FldLocalitynote,omitempty"`
	Start            string  `json:"start,omitempty"`
	End              string  `json:"end,omitempty"`
	State            int     `json:"state,omitempty"`
	Locality         int     `json:"locality,omitempty"`
	Admin            int     `json:"admin,omitempty"`
}

type bakingAndQuantity struct {
	FldBakeryNo         int     `gorm:"column:FldBakeryNo" json:"FldBakeryNo,omitempty"`
	FldBakeryName       string  `gorm:"column:FldBakeryName" json:"FldBakeryName,omitempty"`
	FldBakedQuantity    float32 `gorm:"column:FldBakedQuantity" json:"FldBakedQuantity,omitempty"`
	FldReceivedQuantity float32 `gorm:"column:FldReceivedQuantity" json:"FldReceivedQuantity,omitempty"`
}

func (f flourData) validate() bool {
	if f.Start != "" && f.End != "" && f.State != 0 && f.Locality != 0 && f.Admin != 0 {
		return true
	}
	return false
}

func (f flourData) marshal() []byte {
	d, _ := json.Marshal(&f)
	return d
}
