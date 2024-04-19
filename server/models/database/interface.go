package models

type Interface struct {
	Id        int64  `json:"id"`        //The id of the network interface
	IdMachine int64  `json:"idMachine"` //The id of the machine it is associated with
	Type      string `json:"type"`      //The type of the interface (IPv4 or IPv6)
	IpAddress string `json:"ipAddress"` //The ip address of the interface
	Name      string `json:"name"`      //The name of the interface (from the OS)
}
