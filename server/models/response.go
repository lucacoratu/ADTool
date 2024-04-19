package models

import (
	"encoding/json"
	"io"

	databaseModels "github.com/lucacoratu/ADTool/server/models/database"
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

type AgentsResponse struct {
	Id            int64  `json:"id"`            //The id of the agent
	Name          string `json:"name"`          //The name of the agent
	Username      string `json:"username"`      //The OS username the agent is running as
	DisplayName   string `json:"displayname"`   //The display name of the user the agent is running as
	OsUserId      string `json:"osUserId"`      //The UserId from the OS of the user the agent is running as
	OsUserGroupId string `json:"osUserGroupId"` //The GroupId from the OS of the user the agent is running as
	HomeDirectory string `json:"homeDirectory"` //The home directory of the OS user the agent is running as
}

func (ar *AgentsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ar)
}

type AgentsApiResponse struct {
	Agents []AgentsResponse `json:"agents"` //The list of agents
}

func (aar *AgentsApiResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(aar)
}

type AgentCommandsApiResponse struct {
	Commands []databaseModels.Command `json:"commands"`
}

func (acar *AgentCommandsApiResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(acar)
}
