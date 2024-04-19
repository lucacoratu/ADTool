package models

import (
	"encoding/json"
	"io"
)

// This structure holds information about a network interface found on the machine
type NetworkInterface struct {
	Type    string `json:"type"`    //The type of the interface (IPv4 or IPv6)
	Address string `json:"address"` //The IP address associated with this interface
	Name    string `json:"name"`    //The name of the interface
}

// This structure holds the groups in which the os user that runs the agent is part of
type OsUserGroups struct {
	ID   string `json:"id"`   //The id of the group
	Name string `json:"name"` // The name of the group
}

// This structure holds information about the user the agent is running as on the machine
type OsUser struct {
	Username      string         `json:"username"`      //The username of the user
	DisplayName   string         `json:"displayName"`   //The display name of the user
	UID           string         `json:"uid"`           //The id of the user (on Linux it is an integer, on Windows is a SID (string))
	GID           string         `json:"gid"`           //The group id of the user
	HomeDirectory string         `json:"homeDirectory"` //The home directory of the user
	Groups        []OsUserGroups `json:"groups"`        //The groups the user is part of
}

type MachineInformation struct {
	Hostname      string             `json:"hostname"`
	Os            string             `json:"os"`
	NetInterfaces []NetworkInterface `json:"networkInterfaces"`
	OsCurrentUser OsUser             `json:"osCurrentUser"`
}

func (mi *MachineInformation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mi)
}

func (mi *MachineInformation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mi)
}
