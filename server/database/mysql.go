package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lucacoratu/ADTool/server/configuration"
	"github.com/lucacoratu/ADTool/server/logging"
	"github.com/lucacoratu/ADTool/server/models"
	databaseModels "github.com/lucacoratu/ADTool/server/models/database"
)

// A particular implementation of the IConnection interface
// This is specific for mysql
type MysqlConnection struct {
	logger logging.ILogger
	config configuration.Configuration
	conn   *sql.DB
}

func NewMysqlConnection(logger logging.ILogger, config configuration.Configuration) *MysqlConnection {
	return &MysqlConnection{logger: logger, config: config}
}

func (mysql *MysqlConnection) createTables() error {
	//Create the query to create the machines table
	query := `
		CREATE TABLE IF NOT EXISTS machines(
			id INT PRIMARY KEY AUTO_INCREMENT,
			hostname TEXT,
			os TEXT
		);
	`
	//Execute the query to create the machines table
	_, err := mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the query to create the interfaces table
	query = `
		CREATE TABLE IF NOT EXISTS interfaces (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_machine INT NOT NULL,
			type TEXT, 
			ip_address TEXT,
			name TEXT
		);
	`
	//Execute the query to create the interfaces table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the query to create the agents table
	query = `
		CREATE TABLE IF NOT EXISTS agents (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_machine INT NOT NULL,
			name TEXT,
			username TEXT,
			display_name TEXT, 
			os_user_id TEXT,
			os_user_group_id TEXT,
			home_directory TEXT
		);
	`
	//Execute the query to create the agents table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the query to create the os_groups table
	query = `
		CREATE TABLE IF NOT EXISTS os_groups (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_agent INT NOT NULL,
			os_group_id TEXT,
			os_group_name TEXT
		);
	`
	//Execute the query to create the os_groups table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the table for commands
	query = `
		CREATE TABLE IF NOT EXISTS commands (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_agent INT NOT NULL,
			command TEXT NOT NULL,
			output TEXT
		)
	`
	//Execute the query to create the reccuring commands table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the table for recurring_commands
	query = `
		CREATE TABLE IF NOT EXISTS recurring_commands (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_agent INT NOT NULL,
			command TEXT NOT NULL,
			recurring_interval INT NOT NULL,
			start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	//Execute the query to create the commands table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	//Create the table for recurring_commands_outputs
	query = `
		CREATE TABLE IF NOT EXISTS recurring_commands_outputs (
			id INT PRIMARY KEY AUTO_INCREMENT,
			id_recurring_command INT NOT NULL,
			output TEXT,
			output_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	//Execute the query to create the commands table
	_, err = mysql.conn.Exec(query)
	//Check if an error occured when executing the query
	if err != nil {
		return err
	}

	return nil
}

func (mysql *MysqlConnection) Init() error {
	//Create the connection string
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysql.config.DatabaseUsername, mysql.config.DatabasePassword, mysql.config.DatabaseIPAddress, mysql.config.DatabaseName)
	//Initialize the database connection
	dbConn, err := sql.Open("mysql", connString)
	//Check if an error occured when initializing the connection
	if err != nil {
		return err
	}
	dbConn.SetConnMaxLifetime(time.Minute * 10)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(10)
	//Save the connection instance in the structure
	mysql.conn = dbConn

	//Create the tables of the database if they do not exist
	err = mysql.createTables()

	return err
}

func (mysql *MysqlConnection) RegisterMachine(Hostname string, Os string) (int64, error) {
	//Prepare the query to insert the machine in the database
	query := `
		INSERT INTO machines (hostname, os)
		VALUES (?, ?);
	`
	//Execute the query
	res, err := mysql.conn.Exec(query, Hostname, Os)
	if err != nil {
		return -1, err
	}
	//Get the id of the machine
	id, err := res.LastInsertId()

	return id, err
}

func (mysql *MysqlConnection) RegisterMachineNetworkInterfaces(idMachine int64, netInterfaces []models.NetworkInterface) error {
	for _, netInterface := range netInterfaces {
		//Prepare the query to insert the interface in the database
		query := `
			INSERT INTO interfaces (id_machine, type, ip_address, name)
			VALUES (?,?,?,?)
		`
		_, err := mysql.conn.Exec(query, idMachine, netInterface.Type, netInterface.Address, netInterface.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *MysqlConnection) RegisterAgent(idMachine int64, Username string, DisplayName string, OsUserId string, OsUserGroupId string, HomeDirectory string) (int64, error) {
	//Prepare the query to insert the agent in the database
	query := `
		INSERT INTO agents (id_machine, username, display_name, os_user_id, home_directory)
		VALUES (?,?,?,?,?)
	`
	//Execute the query
	res, err := mysql.conn.Exec(query, idMachine, Username, DisplayName, OsUserId, HomeDirectory)
	if err != nil {
		return -1, err
	}
	agentId, err := res.LastInsertId()
	return agentId, err
}

func (mysql *MysqlConnection) RegisterAgentOSGroups(idAgent int64, groups []models.OsUserGroups) error {
	for _, group := range groups {
		//Prepare the query to insert the group in the database
		query := `
			INSERT INTO os_groups (id_agent, os_group_id, os_group_name)
			VALUES (?,?,?)
		`
		//Execute the query
		_, err := mysql.conn.Exec(query, idAgent, group.ID, group.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mysql *MysqlConnection) RegisterCommand(agentId int64, command string) (int64, error) {
	query := `
		INSERT INTO commands (id_agent, command, output)
		VALUES (?,?,?)
	`
	//Execute the query
	res, err := mysql.conn.Exec(query, agentId, command, "")
	if err != nil {
		return -1, err
	}
	commandId, err := res.LastInsertId()
	return commandId, err
}

func (mysql *MysqlConnection) SetCommandOutput(commandId int64, output string) error {
	query := `
		UPDATE commands SET output = ?
		WHERE id = ?
	`
	//Execute the query
	_, err := mysql.conn.Exec(query, output, commandId)
	return err
}

func (mysql *MysqlConnection) RegisterRecurringCommand(agentId int64, command string, interval int64) (int64, error) {
	query := `
	INSERT INTO recurring_commands (id_agent, command, recurring_interval)
	VALUES (?,?,?)
	`
	//Execute the query
	res, err := mysql.conn.Exec(query, agentId, command, interval)
	if err != nil {
		return -1, err
	}
	commandId, err := res.LastInsertId()
	return commandId, err
}

func (mysql *MysqlConnection) GetAgents() ([]models.AgentsResponse, error) {
	//Prepare the query to get the agents
	query := `
		SELECT id, name, username, display_name, os_user_id, os_user_group_id, home_directory
		FROM agents
	`
	//Execute the query
	rows, err := mysql.conn.Query(query)
	//Check if an error occured when executing the query
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	aux := models.AgentsResponse{}
	returnData := make([]models.AgentsResponse, 0)
	for rows.Next() {
		var name sql.NullString
		var os_user_group_id sql.NullString
		err := rows.Scan(&aux.Id, &name, &aux.Username, &aux.DisplayName, &aux.OsUserId, &os_user_group_id, &aux.HomeDirectory)
		if err != nil {
			return nil, err
		}
		aux.Name = name.String
		aux.OsUserGroupId = os_user_group_id.String
		returnData = append(returnData, aux)
	}
	return returnData, nil
}

func (mysql *MysqlConnection) GetAgentCommands(agentId int64) ([]databaseModels.Command, error) {
	query := `
		SELECT id, command, output
		FROM commands
		WHERE id_agent = ?
		ORDER BY id DESC
	`
	//Execute the query
	rows, err := mysql.conn.Query(query, agentId)
	//Check if an error occured when executing the query
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	aux := databaseModels.Command{}
	returnData := make([]databaseModels.Command, 0)
	for rows.Next() {
		var output sql.NullString
		err := rows.Scan(&aux.Id, &aux.Command, &output)
		if err != nil {
			return nil, err
		}
		aux.Output = output.String
		returnData = append(returnData, aux)
	}
	return returnData, nil
}
