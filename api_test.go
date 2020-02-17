package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_login(t *testing.T) {
	// r := login()

	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	wrongData := Login{Username: "mohamed", Password: "my wrong password"}
	correctData := Login{Username: "admin", Password: "admin"}
	tests := []struct {
		name string
		req  Login
		want int
	}{
		{"400 request", Login{}, 400}, {"200 request", correctData, 200}, {"400 request wrong payload", wrongData, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := marshal(tt.req)
			dd := bytes.NewBuffer(d)
			res, err := http.Post(ts.URL, "application/json", dd)
			if err != nil {
				log.Fatal(err)
			}
			_, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getUser() got = %v, want %v", res.StatusCode, tt.want)
			}
		})
	}
}

func Test_login_contentType(t *testing.T) {
	// r := login()

	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	time1, _ := time.Parse(time.RFC3339, "2019-09-17T18:53:02Z")

	user1 := User{
		FldUserNo:         10,
		FldFullName:       "user doesn't exist",
		FldUserType:       0,
		FldDisabled:       false,
		FldStateNo:        0,
		FldLocaliyNo:      0,
		FldCityNo:         0,
		FldNeighborhoodNo: 0,
		FldSecurityLevel:  0,
		FldUpdateDate:     time1.String(),
	}
	user2 := User{
		/*1,1,Ahmed Mustafa,admin,admin,2
		0x89504E470D0A1A0A0000000D494844520000008C0000008C0806000000AEC0413E000000017352474200AECE1CE90000000467414D410000B18F0BFC6105000000097048597300000EC300000EC301C76FA8640000966A49444154785EA5FD659C5DD7B1F50BF7FBDEE79C24265992C5CCCCCCCCCCCC520BBA052DA9A5965A0x6E85CA3968FA94A862BD4A280CD77D5C52CF92553DB46E21211C5C12093C18A6D913051C309D0B16B5B9EDDADABE49135DC71CD7859C36AC74710FA703BB70E27105C7A78C7237B47FC2103B3A758C5D9A33CDAE2F9C633772E74A8FCCB3AB4B,0,NULL,NULL,NULL,NULL,NULL,2019-09-17 18:53:02.000

		*/
		FldUserNo:         1,
		FldFullName:       "Ahmed Mustafa",
		FldUserType:       0,
		FldDisabled:       false,
		FldStateNo:        0,
		FldLocaliyNo:      0,
		FldCityNo:         0,
		FldNeighborhoodNo: 0,
		FldSecurityLevel:  0,
		FldUpdateDate:     time1.String(),
		FldUserName:       "admin",
	}

	d1 := user1.marshal()
	d2 := user2.marshal()

	wrongData := Login{Username: "mohamed", Password: "my wrong password"}
	correctData := Login{Username: "admin", Password: "admin"}
	tests := []struct {
		name  string
		req   Login
		want  int
		want2 []byte
	}{
		{"400 request", Login{}, 400, d1}, {"200 request", correctData, 200, d2}, {"400 request wrong payload", wrongData, 400, d1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := marshal(tt.req)
			dd := bytes.NewBuffer(d)
			res, err := http.Post(ts.URL, "application/json", dd)
			if err != nil {
				log.Fatal(err)
			}
			got1, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("login handler() got = %v, want %v", res.StatusCode, tt.want)
			} else if res.Header.Get("content-type") != "application/json" {
				t.Errorf("login handler() got = %v, want %v", res.Header.Get("content-type"), "application/json")
			} else if !reflect.DeepEqual(got1, tt.want2) {
				t.Errorf("login handler: got = %v, want = %v", string(got1), tt.want2)
			}
		})
	}
}

func Test_login1(t *testing.T) {
	// r := login()

	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	user1 := User{
		FldUserNo:         10,
		FldFullName:       "user doesn't exist",
		FldUserType:       0,
		FldImage:          nil,
		FldDisabled:       false,
		FldStateNo:        0,
		FldLocaliyNo:      0,
		FldCityNo:         0,
		FldNeighborhoodNo: 0,
		FldSecurityLevel:  0,
		FldUpdateDate:     time.Now().String(),
	}

	time1, _ := time.Parse(time.RFC3339, "2019-09-17T18:53:02Z")

	user2 := User{
		/*1,1,Ahmed Mustafa,admin,admin,2
		0x89504E470D0A1A0A0000000D494844520000008C0000008C0806000000AEC0413E000000017352474200AECE1CE90000000467414D410000B18F0BFC6105000000097048597300000EC300000EC301C76FA8640000966A49444154785EA5FD659C5DD7B1F50BF7FBDEE79C24265992C5CCCCCCCCCCCC520BBA052DA9A5965A0x6E85CA3968FA94A862BD4A280CD77D5C52CF92553DB46E21211C5C12093C18A6D913051C309D0B16B5B9EDDADABE49135DC71CD7859C36AC74710FA703BB70E27105C7A78C7237B47FC2103B3A758C5D9A33CDAE2F9C633772E74A8FCCB3AB4B,0,NULL,NULL,NULL,NULL,NULL,2019-09-17 18:53:02.000

		*/
		FldUserNo:         1,
		FldFullName:       "Ahmed Mustafa",
		FldUserType:       0,
		FldImage:          nil,
		FldDisabled:       false,
		FldStateNo:        0,
		FldLocaliyNo:      0,
		FldCityNo:         0,
		FldNeighborhoodNo: 0,
		FldSecurityLevel:  0,
		FldUpdateDate:     time1.String(),
		FldUserName:       "admin",
	}

	wrongData := Login{Username: "mohamed", Password: "my wrong password"}
	correctData := Login{Username: "admin", Password: "admin"}
	tests := []struct {
		name  string
		req   Login
		want  int
		want2 User
	}{
		{"400 request", Login{}, 400, user1}, {"200 request", correctData, 200, user2}, {"400 request wrong payload", wrongData, 400, user1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := marshal(tt.req)
			dd := bytes.NewBuffer(d)
			res, err := http.Post(ts.URL, "application/json", dd)
			if err != nil {
				log.Fatal(err)
			}
			got1, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			var u User
			json.Unmarshal(got1, &u)

			if res.StatusCode != tt.want {
				t.Errorf("login handler() got = %v, want %v", res.StatusCode, tt.want)
			} else if res.Header.Get("content-type") != "application/json" {
				t.Errorf("login handler() got = %v, want %v", res.Header.Get("content-type"), "application/json")
			} else if !reflect.DeepEqual(u, tt.want2) {
				t.Errorf("login handler: got = %#v, \n\nwant = %#v", u, tt.want2)
			}
		})
	}
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func marshal(l Login) []byte {
	data, _ := json.Marshal(&l)
	return data
}

func Test_getGrinderHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getGrinderHandler))

	defer ts.Close()

	q1 := "agent=2"
	q2 := "agent=3"
	q3 := ""

	g1 := []Grinder{{FldGrinderNo: 3,
		FldGrinderName: "مطاحن سين",
		FldIsActive:    true,
		FldStateNo:     0,
		FldContactName: "N/A",
		FldPhone:       "N/A",
		FldEmail:       "N/A",
		FldAddress:     "N/A",
		FldVolume:      50000,
		FldUserNo:      1,
		FldLogNo:       44,
		FldUpdateDate:  "",
	},
	}
	tests := []struct {
		name  string
		req   string
		want  int
		want2 []Grinder
	}{
		{"Grinder with agent id agent=2", q1, 400, []Grinder{}}, {"grinder with agent id 3", q2, 200, g1},
		{"grinder with agent id 3", q3, 200, g1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL + "?" + tt.req)

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getGrinderFromAgent() got = %v, want %v", res.StatusCode, tt.want)
			}
			ww := marshalGrinder(w)
			if !reflect.DeepEqual(w, tt.want2) {
				t.Errorf("getGrinderFromAgent() got = %v, want %v", ww, tt.want2)

			}
		})
	}
}

func Test_getSubmittedFlourHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getSubmittedFlourHandler))

	defer ts.Close()

	q1 := "agent=2"

	tests := []struct {
		name  string
		req   string
		want  int
		want2 []Grinder
	}{
		{"Grinder with agent id 1", q1, 200, []Grinder{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL)

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getSubmittedFlourHandler() got = %v, want %v\n\nThe body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

// submitFlourHandler
func Test_submitFlourHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(submitFlourHandler))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(_marshalFlourSubmit(tt.args)))

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("submitFloorHandler() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

// test getBakery
func Test_getBakeries(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getBakeries))

	// d := time.Now()
	defer ts.Close()

	tests := []struct {
		name string
		args string
		want int
	}{
		{"case empty request body", "", 200}, {"case empty request body", "", 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL + "?agent=2")

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getBakeries() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

// test submitBakery
func Test_setDistributedFlours(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(setDistributedFlours))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(_marshalFlourSubmit(tt.args)))

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("submitFloorHandler() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

// bakerySubmitFlourHandler
func Test_bakerySubmitFlourHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(bakerySubmitFlourHandler))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(_marshalFlourSubmit(tt.args)))

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("bakerySubmitFlourHandler() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_bakeryAgentsHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(bakeryAgentsHandler))

	defer ts.Close()

	tests := []struct {
		name string
		want int
	}{
		{"case empty request body", 200}, {"case request with all fields", 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL + "?agent=2")

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("bakerySubmitFlourHandler() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

// recordBakedHandler
func Test_recordBakedHandler(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(recordBakedHandler))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(_marshalFlourSubmit(tt.args)))

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("recordBakedHandler() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_listing(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(listing))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL)

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("listing() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_getComplains(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getComplains))

	now := time.Now().UTC()
	d := now.Format("2006-01-02")
	// d := time.Now()
	defer ts.Close()
	qFull := FlourAgentReceive{
		FldFlourAgentReceiveNo:    1,
		FldDate:                   d,
		FldFlourAgentNo:           3,
		FldGrinderNo:              3,
		FldQuantity:               20,
		FldUnitPrice:              3.5,
		FldTotalAmount:            70,
		FldRefNo:                  "",
		FldNFCFlourAgentReceiveNo: 0,
		FldNFCStatusNo:            0,
		FldNFCNote:                "",
		FldUserNo:                 0,
	}
	qMiss := FlourAgentReceive{}

	tests := []struct {
		name string
		args FlourAgentReceive
		want int
	}{
		{"case empty request body", qMiss, 400}, {"case request with all fields", qFull, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL)

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("listing() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_generateComplains(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(generateComplains))

	defer ts.Close()

	tests := []struct {
		name string
		want int
	}{
		{"case empty request body", 400}, {"case request with all fields", 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL)

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("listing() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_getAllBakeries(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAllBakeries))

	// d := time.Now()
	defer ts.Close()

	tests := []struct {
		name string
		args string
		want int
	}{
		{"case empty request body", "", 200}, {"case empty request body", "", 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL + "?agent=2")

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getBakeries() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func Test_getLocations(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getLocations))

	// d := time.Now()
	defer ts.Close()

	tests := []struct {
		name string
		args string
		want int
	}{
		{"case empty request body", "", 200}, {"case empty request body", "", 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Get(ts.URL + "?agent=2&state=1")

			if err != nil {
				log.Fatal(err)
			}

			w, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("getLocations() got = %v, want %v\n\nRes body is: %v", res.StatusCode, tt.want, string(w))
			}
		})
	}
}

func marshalGrinder(d []byte) []Grinder {
	var g []Grinder
	json.Unmarshal(d, &g)
	return g
}

func _marshalFlourSubmit(f FlourAgentReceive) []byte {
	d, _ := json.Marshal(&f)
	return d
}
