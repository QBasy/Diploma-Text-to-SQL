package adapters

import (
	"database/sql"
	"fmt"
)

type MySQLAdapter struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewMySQLAdapter(host string, port int, username, password, database string) *MySQLAdapter {
	return &MySQLAdapter{
		Host: host, Port: port, Username: username,
		Password: password, Database: database,
	}
}

func (a *MySQLAdapter) Connect() (Queryable, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		a.Username, a.Password, a.Host, a.Port, a.Database)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	return conn, nil
}

func (a *MySQLAdapter) GetType() string {
	return "mysql"
}

func (a *MySQLAdapter) Ping() error {
	conn, err := a.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.(*sql.DB).Ping()
}
