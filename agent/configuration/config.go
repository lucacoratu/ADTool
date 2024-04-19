package configuration

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/lucacoratu/ADTool/agent/utils"
)

type Configuration struct {
	ServerURL string `json:"serverURL" validate:"required"` //The URL of the API
	Id        int64  `json:"id"`                            //The id of the agent
}

// Load the configuration from a file
func (conf *Configuration) LoadConfigurationFromFile(filePath string) error {
	//Check if the file exists
	found := utils.CheckFileExists(filePath)
	if !found {
		return errors.New("configuration file cannot be found")
	}
	//Open the file and load the data into the configuration structure
	file, err := os.Open(filePath)
	//Check if an error occured when opening the file
	if err != nil {
		return err
	}
	err = conf.FromJSON(file)
	//Check if an error occured when loading the json from file
	if err != nil {
		return err
	}
	//Initialize the validator of the json data
	validate := validator.New(validator.WithRequiredStructEnabled())
	//Validate the fields of the struct
	err = validate.Struct(conf)
	return err
}

// Convert from json into the configuration structure
func (conf *Configuration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(conf)
}

// Convert to json the configuration structure
func (conf *Configuration) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(conf)
}
