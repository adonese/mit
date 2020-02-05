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
