package models

import (
	"encoding/json"
	"io"
)

type OsGroup struct {
	Id        int64  `json:"id"`        //The id of the group
	IdAgent   int64  `json:"idAgent"`   //The id of the agent it is associated with
	OsGroupId string `json:"osGroupId"` //The id of the group from the OS
	Name      string `json:"name"`      //The name of the group from the OS
}

func (og *OsGroup) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(og)
}

func (og *OsGroup) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(og)
}
