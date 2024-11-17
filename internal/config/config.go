package config // This package is responsible for reading and writing the JSON file

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

/*
	Config represents the JSON file structure:
	{
	"db_url": "connection_string_goes_here",
	"current_user_name": "username_goes_here"
	}
*/

const configFileName = ".gatorconfig.json"

type Config struct {
	dbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read the JSON file located at ~/.gatorconfig.json and return a Config struct.
func Read() (Config, error) {

	filePath, err := getFilePath()
	if err != nil {
		log.Println("Error finding home directory:", err)
		return Config{}, err
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error opening the file:", err)
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config) // Parse the JSON-encoded data and store the result in config
	if err != nil {
		log.Println("Failed to parse data as JSON:", err)
		return Config{}, err
	}

	return config, nil
}

// Set the CurrentUserName field of the Config struct and
// write it to the JSON file
func (config Config) SetUser(user string) error {
	config.CurrentUserName = user
	err := write(config)
	if err != nil {
		log.Println("Error setting current user,", err)
		return err
	}

	return nil
}

// Return the complete fileName of ~/.gatorconfig.json
func getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error finding home directory,", err)
		return "", err
	}

	fileName := filepath.Join(homeDir, configFileName)
	return fileName, nil
}

// Write config to ~/.gatorconfig.json
func write(config Config) error {
	filePath, err := getFilePath()
	if err != nil {
		log.Println("Error finding home directory:", err)
		return err
	}

	bytes, err := json.Marshal(config)
	if err != nil {
		log.Println("Failed to encode data as JSON:", err)
		return err
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return err
	}
	return nil
}
