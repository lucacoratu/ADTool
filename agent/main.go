package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/lucacoratu/ADTool/agent/apiclient"
	"github.com/lucacoratu/ADTool/agent/configuration"
	"github.com/lucacoratu/ADTool/agent/logging"
	"github.com/lucacoratu/ADTool/agent/utils"
	"github.com/lucacoratu/ADTool/agent/websocket"
)

func main() {
	//Initialize the logger
	logger := logging.NewDefaultDebugLogger()
	//Verify the connection to the API via the healthcheck endpoint
	client := &http.Client{}
	_, err := client.Get("http://127.0.0.1:8080/api/v1/healthcheck")
	if err != nil {
		logger.Error("Could not establish the connection to the API", err.Error())
		return
	}

	//Load the configuration from file
	config := configuration.Configuration{}
	err = config.LoadConfigurationFromFile("agent.conf")
	if err != nil {
		logger.Error("Could not load configuration from file", err.Error())
		return
	}

	baseUrl := config.ServerURL + "/api/v1"

	//Register the agent
	if config.Id == 0 {
		//Register the agent if it is not registered already
		machineInfo, err := utils.GetMachineInfo()
		if err != nil {
			logger.Error("Could not get machine information", err.Error())
			return
		}

		apiClient := apiclient.NewAPIClient(logger, baseUrl)
		agentId, err := apiClient.RegisterAgent(machineInfo)
		if err != nil {
			logger.Error("Could not register the agent in the database", err.Error())
		}

		//Save the agent id in the configuration
		logger.Debug("Agent id", agentId)
		config.Id = agentId

		file, err := os.OpenFile("agent.conf", os.O_WRONLY|os.O_TRUNC, 0644)
		//Check if an error occured when trying to open the configuration file to update it
		if err != nil {
			logger.Error("Could not save the configuration file to disk, failed to open configuration file for writing, UUID not saved", err.Error())
		} else {
			newConfigContent, err := json.MarshalIndent(config, "", "    ")
			//Check if an error occured when marshaling the json for configuration
			if err != nil {
				logger.Error("Could not save the configuration file to disk, ID not saved", err.Error())
			} else {
				//Write the new configuration to file
				_, err := file.Write(newConfigContent)
				//Check if an error occured when writing the new configuration
				if err != nil {
					logger.Error("Could not write the new configuration file, ID not saved", err.Error())
				} else {
					logger.Info("Updated the configuration file to contain the received ID from the API")
				}
			}
		}
	}

	//Add the backdoors (SSH public keys in the home directory)

	apiWsConn := websocket.NewAPIWebSocketConnection(logger, "ws://127.0.0.1:8080/api/v1/agents/"+strconv.Itoa(int(config.Id))+"/ws")
	//TO DO... Exponential retry
	_, err = apiWsConn.Connect()
	if err != nil {
		logger.Error("Could not establish the websocket connection to the API", err.Error())
		return
	}
	logger.Info("Agent connected to websocket")
	apiWsConn.Start()
}
