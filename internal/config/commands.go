package config

import (
	"errors"
	"log"
)

// A state holds a pointer to a config
type State struct {
	Ptr *Config
}

// A command contains a name and a slice of string arguments
// gator login <username>
type Command struct {
	Name string   // login
	Args []string // Arguments slice, [<username>] in this case
}

// A commands struct holds all the commands the CLI can handle
type Commands struct {
}

func HandlerLogin(filePath string, s *State, cmd Command) error {
	// Check if cmd.Args contains username
	if len(cmd.Args) == 0 {
		return errors.New("expect username but found none")
	}

	// Set the user to the given username
	username := cmd.Args[0]
	cfg := *(s.Ptr)
	cfg.SetUser(filePath, username)

	log.Printf("Current user is set as %v", username)
	return nil
}
