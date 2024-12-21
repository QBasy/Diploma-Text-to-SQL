package main

import (
	"API/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := handlers.Handlers{}

	auth := router.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	textToSQLApi := router.Group("/text-to-sql")
	{
		textToSQLApi.POST("/convert", handler.ConvertToSQLHandler)
		textToSQLApi.GET("/health", handler.TextToSQLHealthHandler)
	}

	dbApi := router.Group("/database")
	{
		dbApi.POST("/execute-text-to-sql", handler.ExecuteCustomTextToSQLHandler)
		dbApi.POST("/execute-custom-query", handler.ExecuteCustomSQLHandler)

		dbApi.GET("/users", handler.GetUsersHandler)
		dbApi.POST("/users/create", handler.CreateUserHandler)

		dbApi.GET("/get-database", handler.GetUserDatabase)
		dbApi.POST("/create-table", handler.CreateTable)
		dbApi.POST("/execute-query", handler.ExecuteDatabaseQuery)
	}

	customDbApi := router.Group("/db")
	{
		customDbApi.GET("/get-database", handler.GetUserDatabase)
		customDbApi.POST("/create-table", handler.CreateTable)
		customDbApi.POST("/execute-query", handler.ExecuteDatabaseQuery)
	}
	return router
}
