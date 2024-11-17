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
	handlers map[string]func(*State, Command) error // Map of command names to their handler functions
}

func handlerLogin(filePath string, s *State, cmd Command) error {
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

// Register a new handler function for a command name
func (c *Commands) register(name string, f func(*State, Command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*State, Command) error) // Initialize the map if nil
	}
	c.handlers[name] = f
}
