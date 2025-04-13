package routes

import (
	"database-service/controllers"
	"database-service/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	db := connectToPostgreSQL()

	databaseController := controllers.NewDatabaseController(db)

	api := r.Group("/api/database")
	{
		service := api.Group("/")
		{
			service.Use(middleware.VerifyAuthService())
			service.POST("/create-database", databaseController.CreateDatabase)
		}
		api.Use(middleware.AuthMiddleware())
		api.GET("/schema", databaseController.GetDatabaseSchema)
		api.POST("/execute-sql", databaseController.ExecuteSQL)

		api.GET("/schema-complex", databaseController.GetFullDatabaseSchema)

		api.POST("/visualise", databaseController.VisualiseQuery)
	}

	r.GET("/health", controllers.HealthCheck)

	return r
}

func connectToPostgreSQL() *gorm.DB {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	return db
}
