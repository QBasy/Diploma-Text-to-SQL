package main

import (
	"github.com/gin-gonic/gin"
)

func routes(router *gin.Engine) {
	api := router.Group("/user")
	{
		api.POST("/login", loginHandler)
		api.POST("/register", registerHandler)
		api.POST("/query", queryToSQLHandler)
	}

}
