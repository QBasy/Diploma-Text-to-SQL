package main

import (
	"API/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			allowedOrigins := []string{"*"}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := handlers.Handlers{}

	textToSQLApi := router.Group("/text-to-sql")
	{
		textToSQLApi.POST("/convert", handler.ConvertToSQLHandler)
		textToSQLApi.GET("/health", handler.TextToSQLHealthHandler)
	}

	dbApi := router.Group("/database")
	{
		dbApi.POST("/execute-text-to-sql", handler.ExecuteCustomTextToSQLHandler)
		dbApi.POST("/execute-query", handler.ExecuteCustomSQLHandler)
		dbApi.GET("/items", handler.GetItemsHandler)
		dbApi.POST("/items/create", handler.CreateItemHandler)
		dbApi.GET("/users", handler.GetUsersHandler)
		dbApi.POST("/users/create", handler.CreateUserHandler)
	}
	return router
}
