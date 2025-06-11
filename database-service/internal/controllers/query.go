package controllers

import (
	pb "database-service/internal/proto/generated/visualisationpb"
	"database-service/internal/services"
	"database-service/pkg/parsers"
	"database-service/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

type DatabaseController struct {
	db                *gorm.DB
	connectionService *services.DatabaseConnectionService
	queryService      *services.QueryExecutionService
}

func NewDatabaseController(db *gorm.DB) *DatabaseController {
	return &DatabaseController{
		db:                db,
		connectionService: services.NewDatabaseConnectionService(db),
		queryService:      services.NewQueryExecutionService(),
	}
}

func (dc *DatabaseController) VisualiseQuery(c *gin.Context) {
	var request VisualisationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("BIND ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Print(request)

	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		log.Println("USER UUID ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbConn, dbType, err := dc.connectionService.GetConnection(userUUID, request.DatabaseUUID)
	if err != nil {
		log.Println("DB CONNECTION ERROR:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer dbConn.Close()

	queryResult, err := dc.queryService.PrepareForVisualization(dbConn, request.Query)
	if err != nil {
		log.Println("QUERY PREPARATION ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызов сервиса визуализации
	conn, err := grpc.Dial(utils.GetEnv("VISUALISATION_SERVICE", "localhost:5007"), grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to visualisation service"})
		return
	}
	defer conn.Close()

	client := pb.NewVisualisationServiceClient(conn)
	svgResponse, err := client.GenerateChart(c, queryResult)

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

		// Получаем колонки из результата
		cols := make([]string, 0)
		if len(queryResult.Result) > 0 {
			for i := 0; i < len(queryResult.Result[0].Values); i++ {
				cols = append(cols, fmt.Sprintf("column_%d", i+1))
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"svg":                 fallbackSVG,
			"result":              queryResult.Result,
			"columns":             cols,
			"row_count":           len(queryResult.Result),
			"db_type":             dbType,
			"visualization_error": "This query cannot be visualized effectively",
		})
		return
	}

	// Получаем колонки из результата
	cols := make([]string, 0)
	if len(queryResult.Result) > 0 {
		for i := 0; i < len(queryResult.Result[0].Values); i++ {
			cols = append(cols, fmt.Sprintf("column_%d", i+1))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"svg":       svgResponse.Svg,
		"result":    queryResult.Result,
		"columns":   cols,
		"row_count": len(queryResult.Result),
		"db_type":   dbType,
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
		log.Println("BIND ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	dbConn, _, err := dc.connectionService.GetConnection(userUUID, request.DatabaseUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer dbConn.Close()

	queryType, err := dc.queryService.ValidateQuery(request.Query)
	if err != nil {
		log.Println("VALIDATION ERROR")
		log.Println("REQUEST=", request)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if queryType == parsers.QueryTypeSelect {
		result, err := dc.queryService.ExecuteSelectQuery(dbConn, request.Query)
		if err != nil {
			log.Println("EXECUTE SELECT ERROR:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"row_count": result.RowCount,
			"columns":   result.Columns,
			"result":    result.Rows,
		})
		return
	}

	if err := dc.queryService.ExecuteModifyQuery(dbConn, request.Query); err != nil {
		log.Println(request.Query)
		log.Println("EXECUTE MODIFY ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
}
