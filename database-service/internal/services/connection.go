package services

import (
	"database-service/internal/models"
	"database-service/pkg/adapters"
	"database-service/pkg/factories"
	"fmt"
	"gorm.io/gorm"
)

type DatabaseConnectionService struct {
	db      *gorm.DB
	factory *factories.DatabaseAdapterFactory
}

func NewDatabaseConnectionService(db *gorm.DB) *DatabaseConnectionService {
	return &DatabaseConnectionService{
		db:      db,
		factory: factories.NewDatabaseAdapterFactory(),
	}
}

func (s *DatabaseConnectionService) GetConnection(userUUID, dbUUID string) (adapters.Queryable, string, error) {
	if dbUUID == "" {
		return s.getDefaultConnection(userUUID)
	}
	return s.getCustomConnection(userUUID, dbUUID)
}

func (s *DatabaseConnectionService) getDefaultConnection(userUUID string) (adapters.Queryable, string, error) {
	var userDB models.UserDatabase
	if err := s.db.Where("user_uuid = ?", userUUID).First(&userDB).Error; err != nil {
		defaultPath := fmt.Sprintf("./data/users/%s/default.db", userUUID)
		adapter := s.factory.CreateDefaultSQLiteAdapter(defaultPath)
		conn, err := adapter.Connect()
		if err != nil {
			return nil, "", err
		}
		return conn, adapter.GetType(), nil
	}

	adapter := s.factory.CreateDefaultSQLiteAdapter(userDB.Path)
	conn, err := adapter.Connect()
	if err != nil {
		return nil, "", err
	}
	return conn, adapter.GetType(), nil
}

func (s *DatabaseConnectionService) getCustomConnection(userUUID, dbUUID string) (adapters.Queryable, string, error) {
	var customDB models.CustomDatabase
	if err := s.db.Where("uuid = ? AND user_uuid = ?", dbUUID, userUUID).First(&customDB).Error; err != nil {
		return nil, "", fmt.Errorf("database not found or not authorized: %v", err)
	}

	adapter, err := s.factory.CreateAdapter(customDB)
	if err != nil {
		return nil, "", err
	}

	if err := adapter.Ping(); err != nil {
		return nil, "", fmt.Errorf("could not establish connection to %s: %v", adapter.GetType(), err)
	}

	conn, err := adapter.Connect()
	if err != nil {
		return nil, "", err
	}

	return conn, adapter.GetType(), nil
}
