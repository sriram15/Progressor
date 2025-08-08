package service

import (
	"fmt"

	"github.com/sriram15/progressor-todo-app/internal/connection"
)

type ISettingService interface {
	GetAllSettings() (interface{}, error)
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type SettingService struct {
	dbManager *connection.DBManager
	settings  []interface{}
}

func NewSettingService(dbManager *connection.DBManager) *SettingService {
	dbType, dbPath := connection.GetDBInfo()
	settings := []interface{}{
		map[string]string{"key": "dbType", "value": dbType, "display": "Database Type"},
		map[string]string{"key": "dbPath", "value": dbPath, "display": "Database Path"},
		map[string]string{"key": "shortcut_open", "value": "Ctrl + Shift + P", "display": "Shortcut - Open App"},
		map[string]string{"key": "active_card_timeout", "value": "1", "display": "Active Card Timeout (minutes)"},
	}
	return &SettingService{dbManager: dbManager, settings: settings}
}

func (s *SettingService) GetAllSettings() (interface{}, error) {
	return s.settings, nil
}

func (s *SettingService) GetSetting(key string) (string, error) {
	for _, setting_i := range s.settings {
		setting, ok := setting_i.(map[string]string)
		if !ok {
			continue
		}
		if setting["key"] == key {
			return setting["value"], nil
		}
	}
	return "", fmt.Errorf("setting with key '%s' not found", key)
}

func (s *SettingService) SetSetting(key, value string) error {
	// TODO: Implement saving to db
	return nil
}
