package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/ADTool/server/configuration"
	"github.com/lucacoratu/ADTool/server/database"
	"github.com/lucacoratu/ADTool/server/logging"
	"github.com/lucacoratu/ADTool/server/models"
	"github.com/lucacoratu/ADTool/server/websocket"
)

type AgentsHandler struct {
	logger logging.ILogger
	config configuration.Configuration
	dbConn database.IConnection
	wsPool *websocket.Pool
}

func NewAgentsHandler(logger logging.ILogger, config configuration.Configuration, dbConn database.IConnection, wsPool *websocket.Pool) *AgentsHandler {
	return &AgentsHandler{logger: logger, config: config, dbConn: dbConn, wsPool: wsPool}
}

func (ah *AgentsHandler) CreateAgent(rw http.ResponseWriter, r *http.Request) {
	//Create the structure which will hold the data from the agent
	machineInfo := models.MachineInformation{}
	//Get the machine information from the request body
	err := machineInfo.FromJSON(r.Body)
	//Check if an error occured when parsing the request body from the agent
	if err != nil {
		//Send an APIError back to the client
		apiErr := models.NewRequestParseError("Invalid JSON request, check the fields and try again")
		rw.WriteHeader(http.StatusBadRequest)
		apiErr.ToJSON(rw)
		return
	}

	//TO DO.... Validate input provided by the client

	//Register the machine in the database
	machineId, err := ah.dbConn.RegisterMachine(machineInfo.Hostname, machineInfo.Os)
	//Check if an error occured when inserting the machine in the database
	if err != nil {
		//Send an APIError back to the client
		apiErr := models.NewDatabaseError("Could not insert machine in the database")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Register the network interfaces of the machine in the database
	err = ah.dbConn.RegisterMachineNetworkInterfaces(machineId, machineInfo.NetInterfaces)
	//Check if an error occured when inserting the network interfaces for the machine
	if err != nil {
		apiErr := models.NewDatabaseError("Could not insert machine network interfaces in the database")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Register the agent in the database
	agentId, err := ah.dbConn.RegisterAgent(machineId, machineInfo.OsCurrentUser.Username, machineInfo.OsCurrentUser.DisplayName, machineInfo.OsCurrentUser.UID, machineInfo.OsCurrentUser.GID, machineInfo.OsCurrentUser.HomeDirectory)
	//Check if an error occured when inserting the agent in the database
	if err != nil {
		apiErr := models.NewDatabaseError("Could not insert agent in the database")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Register the os groups of the user the agent is running as
	err = ah.dbConn.RegisterAgentOSGroups(agentId, machineInfo.OsCurrentUser.Groups)
	//Check if an error occured
	if err != nil {
		apiErr := models.NewDatabaseError("Could not insert the agent user groups in the database")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Create the return structure
	resp := models.AgentRegisterResponse{Status: "ok", AgentId: agentId}

	//Return the success message
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler to get all the agents registered in the database
func (ah *AgentsHandler) GetAgents(rw http.ResponseWriter, r *http.Request) {
	//Get the agents from the database
	agents, err := ah.dbConn.GetAgents()
	ah.logger.Debug(agents)
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewDatabaseError("Could not get agents")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	resp := models.AgentsApiResponse{Agents: agents}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

func (ah *AgentsHandler) GetCommands(rw http.ResponseWriter, r *http.Request) {
	//Get the agent id from the URL
	vars := mux.Vars(r)
	agent_id, _ := strconv.Atoi(vars["id"])

	//Get the commands from the database
	commands, err := ah.dbConn.GetAgentCommands(int64(agent_id))
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewDatabaseError("Could not get agent's commands")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	resp := models.AgentCommandsApiResponse{Commands: commands}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

func (ah *AgentsHandler) ExecuteCommandOnAgent(rw http.ResponseWriter, r *http.Request) {
	//Get the agent id from the URL
	vars := mux.Vars(r)
	agent_id, _ := strconv.Atoi(vars["id"])

	//Get the command from the body
	cmdMsg := models.ExecuteCommand{}
	err := cmdMsg.FromJSON(r.Body)
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewRequestParseError("Could not parse command from body")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Save the command in the database to get the id
	commandId, err := ah.dbConn.RegisterCommand(int64(agent_id), cmdMsg.Command)
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewDatabaseError("Could not insert the command")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	ah.wsPool.SendExecuteCommandToAgent(int64(agent_id), commandId, cmdMsg.Command)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}

func (ah *AgentsHandler) ExecuteRecurringCommandOnAgent(rw http.ResponseWriter, r *http.Request) {
	//Get the agent id from the URL
	vars := mux.Vars(r)
	agent_id, _ := strconv.Atoi(vars["id"])

	//Get the command from the body
	cmdMsg := models.ExecuteRecurringCommand{}
	err := cmdMsg.FromJSON(r.Body)
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewRequestParseError("Could not parse command from body")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	//Save the command in the database to get the id
	commandId, err := ah.dbConn.RegisterRecurringCommand(int64(agent_id), cmdMsg.Command, cmdMsg.Interval)
	if err != nil {
		ah.logger.Error(err.Error())
		apiErr := models.NewDatabaseError("Could not insert the command")
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr.ToJSON(rw)
		return
	}

	ah.wsPool.SendExecuteRecurringCommandToAgent(int64(agent_id), commandId, cmdMsg.Command, cmdMsg.Interval)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}
