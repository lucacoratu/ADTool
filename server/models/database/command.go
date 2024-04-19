package models

import (
	"encoding/json"
	"io"
)

type Command struct {
	Id      int64  `json:"id"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

func (c *Command) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Command) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}
