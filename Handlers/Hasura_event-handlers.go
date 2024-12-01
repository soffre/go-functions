package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

func Get_hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hellow world",
	})
}

func Comment_trigger(c *gin.Context) {
	var payload struct {
		Event struct {
			Data struct {
				New struct {
					PhotoId string `json:"photo_id"`
				} `json:"new`
			} `json:"data"`
		} `json:"event"`
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
		r.Header.Set("x-hasura-role", "admin")
	})

	err2 := client.Query(context.Background(), &query, variables)
	if err2 != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err2.Error()})
		return
	}

	// if len(query.Photo) == 0{
	// 	c.JSON(http.StatusNotFound, gin.H{ "error": "Not found",})
	// }

	fmt.Println(query.Photo.Photo_url, query.Photo.Name)
	c.JSON(http.StatusOK, gin.H{"message": query.Photo.Photo_url})
	return
}
