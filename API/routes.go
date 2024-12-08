package main

import (
	"github.com/gin-gonic/gin"
)

func routes(router *gin.Engine) {
	textToSQLApi := router.Group("/text-to-sql")
	{
		textToSQLApi.POST("/convert", convertToSQLHandler)
		textToSQLApi.GET("/health", textToSQLHealthHandler)
	}

	dbApi := router.Group("/database")
	{
		dbApi.POST("/execute-query", executeCustomSQLHandler)
		dbApi.GET("/items", getItemsHandler)
		dbApi.POST("/items/create", createItemHandler)
		dbApi.GET("/users", getUsersHandler)
		dbApi.POST("/users/create", createUserHandler)
	}
}
