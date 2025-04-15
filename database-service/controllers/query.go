//query.go file with VisualiseQuery and ExecuteSQL functions
package controllers

import (
	"database-service/models"
	pb "database-service/proto/generated/visualisationpb"
	"database-service/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func (dc *DatabaseController) VisualiseQuery(c *gin.Context) {
	var request VisualisationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(request)
	dbConn, dbTypes, err := GetDatabaseConnection(dc.db, userUUID, request.DatabaseUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query(request.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SQL query"})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch column info"})
		return
	}

	var queryResult pb.QueryResult
	queryResult.SqlQuery = request.Query

	for rows.Next() {
		columnVals := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))
		for i := range columnVals {
			columnPtrs[i] = &columnVals[i]
		}
		if err := rows.Scan(columnPtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process row data"})
			return
		}
		row := &pb.Row{}
		for _, val := range columnVals {
			row.Values = append(row.Values, fmt.Sprintf("%v", val))
		}
		queryResult.Result = append(queryResult.Result, row)
	}

	conn, err := grpc.Dial(utils.GetEnv("VISUALISATION_SERVICE", "localhost:5007"), grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to visualisation service"})
		return
	}
	defer conn.Close()

	client := pb.NewVisualisationServiceClient(conn)
	svgResponse, err := client.GenerateChart(c, &queryResult)

	if err != nil {
		fallbackSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="600" height="200" viewBox="0 0 600 200">
            <rect width="600" height="200" fill="#f8f9fa" rx="10" ry="10" />
            <text x="300" y="80" font-family="Arial, sans-serif" font-size="18" text-anchor="middle" fill="#6c757d">
                We're unable to visualize this query at the moment
            </text>
            <text x="300" y="120" font-family="Arial, sans-serif" font-size="14" text-anchor="middle" fill="#6c757d">
                Try a different query with aggregations or fewer columns
            </text>
        </svg>`

		c.JSON(http.StatusOK, gin.H{
			"svg":                 fallbackSVG,
			"result":              queryResult.Result,
			"columns":             cols,
			"row_count":           len(queryResult.Result),
			"db_type":             dbTypes[0],
			"visualization_error": "This query cannot be visualized effectively",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"svg":       svgResponse.Svg,
		"result":    queryResult.Result,
		"columns":   cols,
		"row_count": len(queryResult.Result),
		"db_type":   dbTypes[0],
	})
}

func (dc *DatabaseController) ExecuteSQL(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request ExecuteSQLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	dbConn, _, err := GetDatabaseConnection(dc.db, userUUID, request.DatabaseUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer dbConn.Close()

	query := strings.TrimSpace(request.Query)
	upperQuery := strings.ToUpper(query)

	switch {
	case isSelectQuery(upperQuery):
		rows, err := dbConn.Query(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SQL query"})
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		var result []map[string]interface{}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPtrs := make([]interface{}, len(cols))
			for i := range columns {
				columnPtrs[i] = &columns[i]
			}
			if err := rows.Scan(columnPtrs...); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
			rowMap := make(map[string]interface{})
			for i, col := range cols {
				rowMap[col] = columns[i]
			}
			result = append(result, rowMap)
		}

		c.JSON(http.StatusOK, gin.H{
			"row_count": len(result),
			"columns":   cols,
			"result":    result,
		})
		return

	case isCreateQuery(upperQuery), isDropQuery(upperQuery), isInsertQuery(upperQuery), isUpdateQuery(upperQuery), isDeleteQuery(upperQuery):
		_, err := dbConn.Exec(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to execute SQL query"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query type not allowed"})
	}
}

func isSelectQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "SELECT")
}

func isCreateQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "CREATE")
}

func isDropQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "DROP")
}

func isInsertQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "INSERT")
}

func isUpdateQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "UPDATE")
}

func isDeleteQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "DELETE")
}

type Queryable interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

func GetDatabaseConnection(db *gorm.DB, userUUID, dbUUID string) (Queryable, []string, error) {
	if dbUUID == "" {
		log.Printf("WTF")
		var userDB models.UserDatabase
		if err := db.Where("user_uuid = ?", userUUID).First(&userDB).Error; err != nil {
			defaultPath := fmt.Sprintf("./data/users/%s/default.db", userUUID)
			sqliteDB, err := sql.Open("sqlite", defaultPath)
			if err != nil {
				return nil, nil, err
			}
			return sqliteDB, []string{"sqlite"}, nil
		}

		sqliteDB, err := sql.Open("sqlite", userDB.Path)
		if err != nil {
			return nil, nil, err
		}
		return sqliteDB, []string{"sqlite"}, nil
	}

	var customDB models.CustomDatabase
	if err := db.Where("uuid = ? AND user_uuid = ?", dbUUID, userUUID).First(&customDB).Error; err != nil {
		return nil, nil, fmt.Errorf("database not found or not authorized: %v", err)
	}

	var dsn string
	switch customDB.DBType {
	case "postgres":
		log.Println("FINE")
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			customDB.Host, customDB.Port, customDB.Username, customDB.Password, customDB.Database, customDB.SSLMode)
		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
		}
		if err := conn.Ping(); err != nil {
			conn.Close()
			return nil, nil, fmt.Errorf("could not establish connection to PostgreSQL: %v", err)
		}
		return conn, []string{"postgres"}, nil

	case "sqlite":
		conn, err := sql.Open("sqlite", customDB.Database)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to SQLite: %v", err)
		}
		if err := conn.Ping(); err != nil {
			conn.Close()
			return nil, nil, fmt.Errorf("could not establish connection to SQLite: %v", err)
		}
		return conn, []string{"sqlite"}, nil

	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			customDB.Username, customDB.Password, customDB.Host, customDB.Port, customDB.Database)
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to MySQL: %v", err)
		}
		if err := conn.Ping(); err != nil {
			conn.Close()
			return nil, nil, fmt.Errorf("could not establish connection to MySQL: %v", err)
		}
		return conn, []string{"mysql"}, nil

	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s", customDB.DBType)
	}
}
