package factories

import (
	"database-service/internal/models"
	"database-service/pkg/adapters"
	"fmt"
)

type DatabaseAdapterFactory struct{}

func NewDatabaseAdapterFactory() *DatabaseAdapterFactory {
	return &DatabaseAdapterFactory{}
}

func (f *DatabaseAdapterFactory) CreateAdapter(customDB models.CustomDatabase) (adapters.DatabaseAdapter, error) {
	switch customDB.DBType {
	case "postgres":
		return adapters.NewPostgreSQLAdapter(
			customDB.Host, customDB.Port, customDB.Username,
			customDB.Password, customDB.Database, customDB.SSLMode,
		), nil
	case "sqlite":
		return adapters.NewSQLiteAdapter(customDB.Database), nil
	case "mysql":
		return adapters.NewMySQLAdapter(
			customDB.Host, customDB.Port, customDB.Username,
			customDB.Password, customDB.Database,
		), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", customDB.DBType)
	}
}

func (f *DatabaseAdapterFactory) CreateDefaultSQLiteAdapter(path string) adapters.DatabaseAdapter {
	return adapters.NewSQLiteAdapter(path)
}
