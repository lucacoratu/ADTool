package websocket

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/lucacoratu/ADTool/server/database"
	"github.com/lucacoratu/ADTool/server/logging"
)

/*
 * This structure will handle concurrent connections using channels
 * Each channel will have a particular functionality
 */
type Pool struct {
	RegisterAgent   chan *AgentClient     //Channel which will handle new agent connections
	UnregisterAgent chan *AgentClient     //Channgel which will handle agent client disconnecting
	AgentClients    map[*AgentClient]bool //A map of dashboard client connections and associated state of the connection (true for online)
	AgentBroadcast  chan AgentMessage     //Channel which will be used to handle a message from the agent
	logger          logging.ILogger       //The logger
	dbConn          database.IConnection  //The database connection
}

/*
 * This function will create a new pool that can then be used when starting the chat service
 */
func NewPool(l logging.ILogger, dbConn database.IConnection) *Pool {
	return &Pool{
		RegisterAgent:   make(chan *AgentClient),
		UnregisterAgent: make(chan *AgentClient),
		AgentClients:    make(map[*AgentClient]bool),
		AgentBroadcast:  make(chan AgentMessage),
		logger:          l,
		dbConn:          dbConn,
	}
}

func (pool *Pool) AgentRegistered(c *AgentClient) {
	c.Status = "online"
	pool.logger.Info("Agent connected to websocket, id:", c.Id)
}

func (pool *Pool) AgentUnregistered(c *AgentClient) {
	c.Status = "offline"
	pool.logger.Info("Agent disconnected from websocket, id: ", c.Id)
}

/*
 * This function will handle when a message is recevied from a client
 * There should be more types of messages that can be received from the client
 */
func (pool *Pool) AgentMessageReceived(message AgentMessage) {
	//Log that a message has been received on the websocket
	pool.logger.Info("Agent message received on the websocket", message.Body)
	//Parse the message body to a websocket message
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
		pool.logger.Debug("Error message received")
	case WsExecuteCommandResponse:
		//Save the command output in the database
		marshaledData, _ := json.Marshal(wsMessage.Data)
		resp := ExecuteCommandResponse{}
		json.Unmarshal(marshaledData, &resp)
		pool.dbConn.SetCommandOutput(resp.Id, resp.Output)
	case WsExecuteRecurringCommandResponse:
		marshaledData, _ := json.Marshal(wsMessage.Data)
		resp := ExecuteCommandResponse{}
		json.Unmarshal(marshaledData, &resp)
		pool.logger.Debug(resp)
	}
}

/*
 * This function will start the pool which will handle client connections, client disconnections and broadcast messages
 */
func (pool *Pool) Start() {
	//Loop infinetly
	for {
		//Check what kind of event occured (connect, disconnect, broadcast message)
		select {
		case client := <-pool.RegisterAgent:
			//Agent connected to the websocket
			pool.AgentClients[client] = true
			pool.logger.Debug("Size of agents connection pool", len(pool.AgentClients))
			pool.AgentRegistered(client)

		case client := <-pool.UnregisterAgent:
			//Agent client disconnected from the websocket
			pool.AgentUnregistered(client)
			delete(pool.AgentClients, client)
			pool.logger.Debug("Size of agents connection pool: ", len(pool.AgentClients))

		case message := <-pool.AgentBroadcast:
			//Message received from the agent on the websocket
			pool.AgentMessageReceived(message)
		}
	}
}

// Function to request the agent to execute a command
func (pool *Pool) SendExecuteCommandToAgent(agentId int64, commandId int64, command string) error {
	for agent := range pool.AgentClients {
		if agent.Id == agentId {
			msg := ExecuteCommandMessage{Id: commandId, Command: command}
			wsMsg := WebSocketMessage{Type: WsExecuteCommand, Data: msg}
			return agent.Conn.WriteJSON(wsMsg)
		}
	}
	return errors.New("agent not found")
}

// Function to request the agent to execute a command
func (pool *Pool) SendExecuteRecurringCommandToAgent(agentId int64, commandId int64, command string, interval int64) error {
	for agent := range pool.AgentClients {
		if agent.Id == agentId {
			msg := ExecuteRecurringCommandMessage{Id: commandId, Command: command, Interval: interval}
			wsMsg := WebSocketMessage{Type: WsExecuteRecurringCommand, Data: msg}
			return agent.Conn.WriteJSON(wsMsg)
		}
	}
	return errors.New("agent not found")
}
