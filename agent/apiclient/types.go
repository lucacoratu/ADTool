package apiclient

import (
	"encoding/json"
	"io"
)

type AgentRegisterResponse struct {
	Status  string `json:"status"`
	AgentId int64  `json:"agentId"`
}

func (arr *AgentRegisterResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(arr)
}

func (arr *AgentRegisterResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(arr)
}
