package models

import (
	"encoding/json"
	"io"
)

type ExecuteCommand struct {
	Command string `json:"command"`
}

func (ec *ExecuteCommand) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ec)
}

type ExecuteRecurringCommand struct {
	Command  string `json:"command"`
	Interval int64  `json:"interval"`
}

func (erc *ExecuteRecurringCommand) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(erc)
}
