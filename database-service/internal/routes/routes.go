package routes

import (
	"database-service/internal/controllers"
	middleware2 "database-service/internal/middleware"
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
			service.Use(middleware2.VerifyAuthService())
			service.POST("/create-database", databaseController.CreateDatabase)
		}
		api.Use(middleware2.AuthMiddleware())
		api.GET("/schema", databaseController.GetDatabaseSchema)
		api.POST("/execute-sql", databaseController.ExecuteSQL) // заменяем на новый контроллер
		api.GET("/schema-complex", databaseController.GetFullDatabaseSchema)
		api.POST("/visualise", databaseController.VisualiseQuery) // заменяем на новый контроллер

		custom := api.Group("/custom")
		{
			custom.POST("/add", databaseController.AddCustomDatabase)
			custom.DELETE("/delete", databaseController.DeleteCustomDatabase)

			custom.GET("/schema", databaseController.GetCustomDatabaseSchema)
			custom.GET("/schema-complex", databaseController.GetFullCustomDatabaseSchema)

			custom.GET("/list", databaseController.ListCustomDatabases)
		}
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
