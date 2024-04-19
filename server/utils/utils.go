package utils

import "os"

// Check if the filepath is valid and exists on the disk
func CheckFileExists(filePath string) bool {
	//Get the current directory
	//pwd, _ := os.Getwd()
	//Check if the file exists
	_, err := os.Stat(filePath)
	//Return the result
	return !os.IsNotExist(err)
}
