package controllers

import (
	"database-service/models"
	pb "database-service/proto/generated/visualisationpb"
	"database-service/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	var userDatabase models.UserDatabase
	if err := dc.db.Where("user_uuid = ?", userUUID).First(&userDatabase).Error; err != nil {
		log.Printf("Error fetching database for UUID %s: %v", userUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer sqliteDB.Close()
	log.Printf("1 %v", request.Query)
	rows, err := sqliteDB.Query(request.Query)
	if err != nil {
		log.Printf("Error querying database for query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SQL query"})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch column info"})
		return
	}

	log.Printf("cols: %v", cols)
	var queryResult pb.QueryResult
	queryResult.SqlQuery = request.Query

	for rows.Next() {
		columnVals := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))
		for i := range columnVals {
			columnPtrs[i] = &columnVals[i]
		}

		if err := rows.Scan(columnPtrs...); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process row data"})
			return
		}

		rowData := pb.Row{}
		for _, col := range columnVals {
			rowData.Values = append(rowData.Values, fmt.Sprintf("%v", col))
		}
		queryResult.Result = append(queryResult.Result, &rowData)
	}

	log.Println("1")
	conn, err := grpc.Dial(utils.GetEnv("VISUALISATION_SERVICE", "localhost:5007"), grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to visualisation service"})
		return
	}
	defer conn.Close()
	log.Println("2")
	client := pb.NewVisualisationServiceClient(conn)
	svgResponse, err := client.GenerateChart(c, &queryResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate visualization"})
		return
	}
	log.Println("3")
	c.JSON(http.StatusOK, gin.H{
		"svg":       svgResponse.Svg,
		"result":    queryResult.Result,
		"columns":   cols,
		"row_count": len(queryResult.Result),
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

	var userDatabase models.UserDatabase
	if err := dc.db.Where("user_uuid = ?", userUUID).First(&userDatabase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer sqliteDB.Close()

	if !(isSelectQuery(request.Query) || isCreateQuery(request.Query) || isDropQuery(request.Query) || isInsertQuery(request.Query) || isUpdateQuery(request.Query) || isDeleteQuery(request.Query)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query type not allowed"})
		return
	}

	if isSelectQuery(request.Query) {
		rows, err := sqliteDB.Query(request.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SQL query"})
			return
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch column information"})
			return
		}

		var result []map[string]interface{}
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPtrs := make([]interface{}, len(cols))
			for i := range columns {
				columnPtrs[i] = &columns[i]
			}

			if err := rows.Scan(columnPtrs...); err != nil {
				log.Printf("Error scanning row: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process row data"})
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
	}

	_, err = sqliteDB.Exec(request.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to execute SQL query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
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
