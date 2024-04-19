package websocket

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lucacoratu/ADTool/agent/logging"
)

type message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type APIWebSocketConnection struct {
	logger     logging.ILogger //The logger
	apiWsURL   string          //The ws url of the API
	State      bool            //The state of the websocket connection (true for active, false for inactive)
	connection *websocket.Conn //The connection structure
}

func NewAPIWebSocketConnection(logger logging.ILogger, apiWsURL string) *APIWebSocketConnection {
	return &APIWebSocketConnection{logger: logger, apiWsURL: apiWsURL}
}

// Connects to the API websocket URL for the agent
func (awsc *APIWebSocketConnection) Connect() (bool, error) {
	//Connect to the websocket URL from the API
	c, _, err := websocket.DefaultDialer.Dial(awsc.apiWsURL, nil)

	//Check if an error occured
	if err != nil {
		return false, err
	}

	awsc.connection = c
	awsc.State = true
	return true, nil
}

// // Handle the connection closed
// func (awsc *APIWebSocketConnection) connectionClosed(code int, text string) error {
// 	return nil
// }

// Handle the message received
func (awsc *APIWebSocketConnection) handleReceivedMessage(message message) {
	awsc.logger.Debug("Message received", message)
	wsMessage := WebSocketMessage{}
	err := wsMessage.FromJSON(strings.NewReader(message.Body))
	//Check if an error occured when parsing the WebSocketMessage from JSON
	if err != nil {
		//Send an error message back to the client
		// errMessage := WebSocketMessage{Type: WsError, Data: data.APIError{Code: data.PARSE_ERROR, Message: "Cannot parse the websocket message from JSON"}}
		// message.C.Conn.WriteJSON(errMessage)
		return
	}
	//Select the action based on the message type
	switch wsMessage.Type {
	case WsError:
		awsc.logger.Debug("Error message received")
	case WsExecuteCommand:
		awsc.logger.Debug("Execute system command")
		id, output, _ := ExecuteSystemCommand(wsMessage)
		//Check if an error occured
		resp := ExecuteCommandResponse{Id: id, Output: output}
		wsRespMsg := WebSocketMessage{Type: WsExecuteCommandResponse, Data: resp}
		//Send the response back to the api
		awsc.connection.WriteJSON(wsRespMsg)
	case WsExecuteRecurringCommand:
		awsc.logger.Debug("Execute recurring system command")
		ExecuteRecurringSystemCommand(wsMessage, awsc.connection)
	}
}

func (awsc *APIWebSocketConnection) Start() {
	//Close the connection at the end of the function
	defer awsc.connection.Close()

	//Start listening for incomming messages
	for {
		mt, msg, err := awsc.connection.ReadMessage()
		if err != nil {
			//Wait a bit then retry the connection
			awsc.logger.Error(err.Error())
			time.Sleep(time.Second * 10)
			_, err = awsc.Connect()
			if err == nil {
				//Connection was restored
				awsc.logger.Info("WebSocket connection to the API has been restored")
			}
			continue
		}

		//Call the handle message function
		awsc.handleReceivedMessage(message{Type: mt, Body: string(msg)})
	}
}
