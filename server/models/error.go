package models

import (
	"encoding/json"
	"io"
)

// This structure holds the information about the error message that the api can send back to the clients
type APIError struct {
	Code    int64  `json:"code"`    //The code of the error
	Message string `json:"message"` //The message of the error
}

func (ae *APIError) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ae)
}

func (ae *APIError) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ae)
}

// Constant values for the codes of the APIErrors
const (
	APIRequestParseError int64 = 0
	DatabaseError        int64 = 1
)

func NewRequestParseError(message string) APIError {
	return APIError{Code: APIRequestParseError, Message: message}
}

func NewDatabaseError(message string) APIError {
	return APIError{Code: DatabaseError, Message: message}
}
