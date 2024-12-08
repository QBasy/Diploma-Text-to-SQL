package main

import (
	"github.com/gin-gonic/gin"
)

var logger *gin.LoggerConfig

func main() {
	logger = gin.Logger()

	r := gin.Default()

	routes(r)

	_ = r.Run(":5001")
}
