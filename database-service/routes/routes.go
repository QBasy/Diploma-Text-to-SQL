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

	api := r.Group("/api")
	{
		api.Use(middleware.AuthMiddleware())
		api.POST("/create-database", databaseController.CreateDatabase)
		api.POST("/execute-sql", databaseController.ExecuteSQL)
	}

	r.GET("/health", controllers.HealthCheck)

	return r
}

func connectToPostgreSQL() *gorm.DB {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	return db
}
