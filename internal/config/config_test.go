package config_test

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jessie-sr/rss-aggregator/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, ".gatorconfig.json")

	mockConfig := config.Config{
		DBUrl:           "mock_db_url",
		CurrentUserName: "mock_user",
	}
	mockData, _ := json.Marshal(mockConfig)
	err := os.WriteFile(tempFile, mockData, 0644)
	assert.NoError(t, err)

	cfg, err := config.Read(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, "mock_db_url", cfg.DBUrl)
	assert.Equal(t, "mock_user", cfg.CurrentUserName)
}

func TestSetUser(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, ".gatorconfig.json")

	mockConfig := config.Config{
		DBUrl:           "mock_db_url",
		CurrentUserName: "mock_user",
	}
	mockData, _ := json.Marshal(mockConfig)
	err := os.WriteFile(tempFile, mockData, 0644)
	assert.NoError(t, err)

	cfg, err := config.Read(tempFile)
	assert.NoError(t, err)

	err = cfg.SetUser(tempFile, "new_mock_user")
	assert.NoError(t, err)

	updatedCfg, err := config.Read(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, "new_mock_user", updatedCfg.CurrentUserName)
}

func TestHandlerLogin_NoArgs(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, ".gatorconfig.json")

	state := &config.State{
		Ptr: &config.Config{},
	}
	cmd := config.Command{
		Name: "login",
		Args: []string{}, // No username provided
	}

	err := config.HandlerLogin(tempFile, state, cmd)
	assert.EqualError(t, err, "expect username but found none")
}

func TestHandlerLogin_ValidUsername(t *testing.T) {
	// Create a temporary file to mock ~/.gatorconfig.json
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, ".gatorconfig.json")

	// Mock state and command
	state := &config.State{
		Ptr: &config.Config{},
	}
	cmd := config.Command{
		Name: "login",
		Args: []string{"testuser"},
	}

	// Capture log output
	var logs strings.Builder
	log.SetOutput(&logs)
	defer log.SetOutput(os.Stderr)

	err := config.HandlerLogin(tempFile, state, cmd)
	assert.NoError(t, err)

	// Verify log contains the correct message
	assert.Contains(t, logs.String(), "Current user is set as testuser")

	// Verify file content
	data, err := config.Read(tempFile)
	assert.NoError(t, err)

	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"current_user_name":"testuser"`)
}
