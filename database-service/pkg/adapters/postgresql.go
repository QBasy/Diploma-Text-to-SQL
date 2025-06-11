package adapters

import (
	"database/sql"
	"fmt"
)

type PostgreSQLAdapter struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewPostgreSQLAdapter(host string, port int, username, password, database, sslmode string) *PostgreSQLAdapter {
	return &PostgreSQLAdapter{
		Host: host, Port: port, Username: username,
		Password: password, Database: database, SSLMode: sslmode,
	}
}

func (a *PostgreSQLAdapter) Connect() (Queryable, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		a.Host, a.Port, a.Username, a.Password, a.Database, a.SSLMode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	return conn, nil
}

func (a *PostgreSQLAdapter) GetType() string {
	return "postgres"
}

func (a *PostgreSQLAdapter) Ping() error {
	conn, err := a.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.(*sql.DB).Ping()
}
