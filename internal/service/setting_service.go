package service

import (
	"github.com/sriram15/progressor-todo-app/internal/connection"
)

type SettingService interface {
	GetAllSettings() (interface{}, error)
	SetSetting(key, value string) error
}

type settingService struct {
}

func NewSettingService() SettingService {
	return &settingService{}
}

func (s *settingService) GetAllSettings() (interface{}, error) {

	dbType, dbPath := connection.GetDBInfo()
	settings := []interface{}{
		map[string]string{"key": "dbType", "value": dbType, "display": "Database Type"},
		map[string]string{"key": "dbPath", "value": dbPath, "display": "Database Path"},
		map[string]string{"key": "shortcut_open", "value": "Ctrl + Shift + P", "display": "Shortcut - Open App"},
	}

	return settings, nil
}

func (s *settingService) SetSetting(key, value string) error {
	// TODO: Implement saving to db
	return nil
}
