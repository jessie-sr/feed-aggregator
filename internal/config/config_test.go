package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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
