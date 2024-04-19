package database

import (
	"github.com/lucacoratu/ADTool/server/models"
	databaseModels "github.com/lucacoratu/ADTool/server/models/database"
)

// This structure is the interface for interacting with database
// It contains all the functions needed by the server
type IConnection interface {
	Init() error
	RegisterMachine(Hostname string, Os string) (int64, error)
	RegisterMachineNetworkInterfaces(idMachine int64, netInterfaces []models.NetworkInterface) error
	RegisterAgent(idMachine int64, Username string, DisplayName string, OsUserId string, osUserGroupId string, HomeDirectory string) (int64, error)
	RegisterAgentOSGroups(idAgent int64, groups []models.OsUserGroups) error
	RegisterCommand(agentId int64, command string) (int64, error)
	RegisterRecurringCommand(agentId int64, command string, interval int64) (int64, error)
	SetCommandOutput(commandId int64, output string) error
	GetAgents() ([]models.AgentsResponse, error)
	GetAgentCommands(agentId int64) ([]databaseModels.Command, error)
}
