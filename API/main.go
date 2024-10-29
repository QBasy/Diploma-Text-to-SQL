package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes(r)

	_ = r.Run(":5001")
}
