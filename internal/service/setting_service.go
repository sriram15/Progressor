package service

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/sriram15/progressor-todo-app/internal/connection"
)

type ISettingService interface {
	GetAllSettings() ([]SettingsItem, error)
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type SettingService struct {
	dbManager *connection.DBManager
	settings  []SettingsItem
}

type SettingsItem struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Display string `json:"display"`
}

func NewSettingService(dbManager *connection.DBManager) *SettingService {
	dbType, dbPath := connection.GetDBInfo()
	settings := []SettingsItem{
		{Key: "dbType", Value: dbType, Display: "Database Type" },
		{Key: "dbPath", Value: dbPath, Display: "Database Path"},
		// map[string]string{"key": "shortcut_open", "value": "Ctrl + Shift + P", "display": "Shortcut - Open App"},
		// map[string]string{"key": "active_card_timeout", "value": "1", "display": "Active Card Timeout (minutes)"},
	}

	log.Printf("SettingService initialized with settings: %v", settings)

	return &SettingService{dbManager: dbManager, settings: settings}
}

func (s *SettingService) GetAllSettings() ([]SettingsItem, error) {
	return s.settings, nil
}

func (s *SettingService) GetSetting(key string) (string, error) {
	for _, setting_i := range s.settings {
		if setting_i.Key == key {
			return setting_i.Value, nil
		}
	}
	return "", fmt.Errorf("setting with key '%s' not found", key)
}

func (s *SettingService) SetSetting(key, value string) error {
	// TODO: Implement saving to db
	return nil
}
