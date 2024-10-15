package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	r.POST("/name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, gin.H{

			"message": fmt.Sprintf("hello, %s!", name),
		})
	})

	r.Run(fmt.Sprintf(":%s", port))
}
