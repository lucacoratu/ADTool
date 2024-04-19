package websocket

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// Execute a system command and returns the output
// If an error occured then the error is returned
func ExecuteSystemCommand(msg WebSocketMessage) (int64, string, error) {
	data, _ := json.Marshal(msg.Data)
	cmdMessage := ExecuteCommandMessage{}
	_ = json.Unmarshal(data, &cmdMessage)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd.exe", "/c", cmdMessage.Command)
		output, err := cmd.Output()
		return cmdMessage.Id, string(output), err
	}
	cmd := exec.Command("bash", "-c", cmdMessage.Command)
	output, err := cmd.Output()
	return cmdMessage.Id, string(output), err
}

// Execute a system command every x seconds
func ExecuteRecurringSystemCommand(msg WebSocketMessage, wsConn *websocket.Conn) error {
	data, _ := json.Marshal(msg.Data)
	cmdMessage := ExecuteRecurringCommandMessage{}
	_ = json.Unmarshal(data, &cmdMessage)

	ticker := time.NewTicker(time.Second * time.Duration(cmdMessage.Interval))
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var output string = ""
				if runtime.GOOS == "windows" {
					cmd := exec.Command("cmd.exe", "/c", cmdMessage.Command)
					cmdOut, err := cmd.Output()
					if err != nil {
						continue
					}
					output = string(cmdOut)
				} else {
					cmd := exec.Command("bash", "-c", cmdMessage.Command)
					cmdOut, err := cmd.Output()
					if err != nil {
						continue
					}
					output = string(cmdOut)
				}
				resp := ExecuteCommandResponse{Id: cmdMessage.Id, Output: output}
				wsRespMsg := WebSocketMessage{Type: WsExecuteRecurringCommandResponse, Data: resp}
				//Send the output to the api
				wsConn.WriteJSON(wsRespMsg)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}
