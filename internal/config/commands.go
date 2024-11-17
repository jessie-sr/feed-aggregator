package config

import (
	"errors"
	"log"
)

// A state holds a pointer to a config
type state struct {
	Ptr *Config
}

// A command contains a name and a slice of string arguments
// gator login <username>
type command struct {
	Name string   // login
	Args []string // Arguments slice, [<username>] in this case
}

// A commands struct holds all the commands the CLI can handle
type commands struct {
}

func handlerLogin(s *state, cmd command) error {
	// Check if cmd.Args contains username
	if len(cmd.Args) == 0 {
		return errors.New("expect username but found none")
	}

	filePath, err := GetFilePath()
	if err != nil {
		log.Println("Error finding the file:", err)
		return err
	}

	// Set the user to the given username
	username := cmd.Args[0]
	cfg := *(s.Ptr)
	cfg.SetUser(filePath, username)

	log.Printf("Current user is set as %v", username)
	return nil
}
