package adapters

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type SQLiteAdapter struct {
	Path string
}

func NewSQLiteAdapter(path string) *SQLiteAdapter {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to get absolute path: %v\n", err)
		absPath = path
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("WARNING: SQLite database file does not exist: %s\n", absPath)
	}

	return &SQLiteAdapter{Path: absPath}
}

func (a *SQLiteAdapter) Connect() (Queryable, error) {
	conn, err := sql.Open("sqlite", a.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %v", err)
	}
	return conn, nil
}

func (a *SQLiteAdapter) GetType() string {
	return "sqlite"
}

func (a *SQLiteAdapter) Ping() error {
	conn, err := a.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.(*sql.DB).Ping()
}
