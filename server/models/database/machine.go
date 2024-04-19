package models

import (
	"encoding/json"
	"io"
)

type Machine struct {
	Id       int64  `json:"id"`       //The id of the machine
	Hostname string `json:"hostname"` //The hostname of the machine
	Os       string `json:"os"`       //The operating system of the machine
}

func (m *Machine) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(m)
}

func (m *Machine) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}
