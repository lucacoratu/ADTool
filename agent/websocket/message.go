package websocket

import (
	"encoding/json"
	"io"
)

// Message Types
const (
	WsError                           int64 = -1 //Error message
	WsExecuteCommand                  int64 = 1  //Execute system command message
	WsExecuteCommandResponse          int64 = 2  //Response for execute system command message
	WsExecuteRecurringCommand         int64 = 3  //Execute recurring system command
	WsExecuteRecurringCommandResponse int64 = 4  //Response for execute recurring system command
)

// WebSocket message format
type WebSocketMessage struct {
	Type int64       `json:"type"` //The type of the message
	Data interface{} `json:"data"` //The data of the message as interface (can be any struct)
}

func (wsm *WebSocketMessage) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(wsm)
}

func (wsm *WebSocketMessage) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(wsm)
}

type ExecuteCommandMessage struct {
	Id      int64  `json:"id"`
	Command string `json:"command"`
}

type ExecuteCommandResponse struct {
	Id     int64  `json:"id"`
	Output string `json:"output"`
}

func (ecr *ExecuteCommandResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ecr)
}

type ExecuteRecurringCommandMessage struct {
	Id       int64  `json:"id"`
	Command  string `json:"command"`
	Interval int64  `json:"interval"`
}
