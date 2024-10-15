package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
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

		client := graphql.NewClient("http://localhost:8081/v1/graphql", nil)
		client = client.WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-hasura-admin-secret", "myhasura")
		})
		c.JSON(http.StatusOK, gin.H{
			"message": "nothing",
		})

	})

	r.Run(fmt.Sprintf(":%s", port))
}
