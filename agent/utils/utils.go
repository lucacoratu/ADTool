package utils

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/lucacoratu/ADTool/agent/models"
)

// Check if the filepath is valid and exists on the disk
func CheckFileExists(filePath string) bool {
	//Get the current directory
	//pwd, _ := os.Getwd()
	//Check if the file exists
	_, err := os.Stat(filePath)
	//Return the result
	return !os.IsNotExist(err)
}

// Read all data from a file in a single string
func ReadAllDataFromFile(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	return string(fileData), err
}

// Read all lines in the file
func ReadLinesFromFile(filePath string) ([]string, error) {
	//Get the current directory
	//pwd, _ := os.Getwd()
	//Check if the file exists
	exists := CheckFileExists(filePath)
	if !exists {
		return nil, errors.New("file does not exist")
	}
	//Open the file
	file, err := os.Open(filePath)
	//Check if an error occured when opening the file
	if err != nil {
		return nil, err
	}
	//Close the file at the end of the function
	defer file.Close()

	//Read lines from the file and append it to returning splice
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	lines := []string{}
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	//Return the lines
	return lines, nil
}

// Check connection to the collector
func CheckAPIConnection(apiBaseURL string) bool {
	response, err := http.Get(apiBaseURL + "/healthcheck")
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// Get the current user which is running the application
func GetCurrentUser() (models.OsUser, error) {
	currentUser, err := user.Current()
	if err != nil {
		return models.OsUser{}, err
	}
	groupIds, err := currentUser.GroupIds()
	if err != nil {
		return models.OsUser{}, err
	}

	osUser := models.OsUser{Username: currentUser.Username, DisplayName: currentUser.Name, UID: currentUser.Uid, GID: currentUser.Gid, HomeDirectory: currentUser.HomeDir, Groups: make([]models.OsUserGroups, 0)}
	for _, groupId := range groupIds {
		group, err := user.LookupGroupId(groupId)
		if err != nil {
			fmt.Println("Could not find the group with id", groupId)
			continue
		}
		osUser.Groups = append(osUser.Groups, models.OsUserGroups{ID: groupId, Name: group.Name})
	}

	return osUser, nil
}

// Collects information about the machine
func GetMachineInfo() (models.MachineInformation, error) {
	machineInfo := models.MachineInformation{}
	//Get the operating system
	machineInfo.Os = runtime.GOOS
	//Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		return machineInfo, errors.New("cannot get the hostname of the machine, " + err.Error())
	}
	machineInfo.Hostname = hostname
	//Get the ip addresses on all network interfaces
	ifaces, err := net.Interfaces()
	//Check if an error occured when getting the network interfaces of the machine
	if err != nil {
		return machineInfo, errors.New("cannot get the network interfaces of the machine, " + err.Error())
	}

	//Go through all the network interfaces and add
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			var ipType string
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				//Check if the IP address is IPv4 or IPv6
				if ip.To4() != nil {
					ipType = "ipv4"
				} else {
					ipType = "ipv6"
				}
			case *net.IPAddr:
				ip = v.IP
				//Check if IP address is IPv4 or IPv6
				if ip.To4() != nil {
					ipType = "ipv4"
				} else {
					ipType = "ipv6"
				}
			}
			// process IP address if it is not loopback
			if !ip.IsLoopback() {
				machineInfo.NetInterfaces = append(machineInfo.NetInterfaces, models.NetworkInterface{Type: ipType, Address: ip.String(), Name: i.Name})
			}
		}
	}

	//Get the current user
	user, err := GetCurrentUser()
	//Check if an error occured when getting the current user
	if err != nil {
		return machineInfo, err
	}
	machineInfo.OsCurrentUser = user

	return machineInfo, nil
}
