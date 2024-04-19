package models

import (
	"encoding/json"
	"io"
)

type Agent struct {
	Id            int64  `json:"id"`            //The id of the agent
	IdMachine     int64  `json:"idMachine"`     //The id of the machine the agent is deployed on
	Name          string `json:"name"`          //The name of the agent
	Username      string `json:"username"`      //The OS username the agent is running as
	DisplayName   string `json:"displayname"`   //The display name of the user the agent is running as
	OsUserId      string `json:"osUserId"`      //The UserId from the OS of the user the agent is running as
	OsUserGroupId string `json:"osUserGroupId"` //The GroupId from the OS of the user the agent is running as
	HomeDirectory string `json:"homeDirectory"` //The home directory of the OS user the agent is running as
}

func (a *Agent) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a *Agent) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}
