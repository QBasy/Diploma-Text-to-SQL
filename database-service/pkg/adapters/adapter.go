package adapters

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type Queryable interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

type DatabaseAdapter interface {
	Connect() (Queryable, error)
	GetType() string
	Ping() error
}
