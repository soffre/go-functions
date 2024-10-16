package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		var payload struct {
			Event struct {
				Data struct {
					New struct {
						PhotoId uuid.UUID `json: "photo_id"`
					} `json: "new`
				} `json: "data"`
			} `json: "event"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var query struct {
			Photo struct {
				Photo_url string `graphql:"photo_url`
				Name      string `graphql:"name"`
			} `graphql:"photos_by_pk(id:$id)"`
		}

		variables := map[string]interface{}{
			"id": payload.Event.Data.New.PhotoId,
		}
		client := graphql.NewClient("http://localhost:8081/v1/graphql", nil)
		client = client.WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-hasura-admin-secret", "myhasura")
		})

		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		// if len(query.Photo) == 0{
		// 	c.JSON(http.StatusNotFound, gin.H{ "error": "Not found",})
		// }

		fmt.Println(query.Photo.Photo_url, query.Photo.Name)
		c.JSON(http.StatusOK, gin.H{"message": query.Photo.Photo_url})
		return
	})

	r.Run(fmt.Sprintf(":%s", port))
}
