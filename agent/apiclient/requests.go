package apiclient

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lucacoratu/ADTool/agent/logging"
	"github.com/lucacoratu/ADTool/agent/models"
)

type APIClient struct {
	baseURL string
	logger  logging.ILogger
}

func NewAPIClient(logger logging.ILogger, baseURL string) *APIClient {
	return &APIClient{logger: logger, baseURL: baseURL}
}

// Register a new agent to the api
func (ac *APIClient) RegisterAgent(machineInfo models.MachineInformation) (int64, error) {
	//Send the request to the server to register the agent in the database
	client := http.Client{}
	url := ac.baseURL + "/agents"
	marshaledData, err := json.Marshal(machineInfo)
	//Check if an error occured
	if err != nil {
		return -1, err
	}
	response, err := client.Post(url, "application/json", strings.NewReader(string(marshaledData)))
	if err != nil {
		return -1, err
	}

	resp := AgentRegisterResponse{}
	err = resp.FromJSON(response.Body)
	//Check if an error occured when parsing the response from the server
	if err != nil {
		return -1, err
	}

	return resp.AgentId, nil
}
