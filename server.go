package main

import (
	"fmt"
	"hasura/go-functions/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	routes.Setup(r)
	r.Run(fmt.Sprintf(":%s", port))
}
