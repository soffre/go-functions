package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            map[string]interface{} `json:"input"`
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func Create_user(c *gin.Context) {

	var actionPayload ActionPayload

	if err := c.ShouldBindJSON(&actionPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the credentials map from the Input
	credentials, ok := actionPayload.Input["credentials"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credentials not provided or invalid"})
		return
	}
	// Extract the password map from the credentials
	password, ok := credentials["password"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password not provided or invalid"})
		return
	}
	hashedPassword := HashPassword(password)
	credentials["password"] = hashedPassword

	// Define the mutation structure
	// var mutation struct {
	// 	InsertUsersOne struct {
	// 		ID int `graphql:"Id"`
	// 	} `graphql:"insert_users(email:$email, name:$name, password:$password)"`
	// }
	var mutation struct {
		InsertUsersOne struct {
			ID    string `graphql:"id"`
			Email string `graphql:"email"`
			Name  string `graphql:"name"`
		} `graphql:"insert_users_one(object: {email: $email, name: $name, password: $password})"`
	}
	// Define the variables
	variables := map[string]interface{}{
		"email":    credentials["email"].(string),
		"name":     credentials["name"].(string),
		"password": credentials["password"].(string),
	}

	client := graphql.NewClient("http://localhost:8081/v1/graphql", nil)
	client = client.WithRequestModifier(func(r *http.Request) {
		r.Header.Set("x-hasura-admin-secret", "myhasura")
	})
	err2 := client.Mutate(context.Background(), &mutation, variables)
	if err2 != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User": gin.H{
			"id":    mutation.InsertUsersOne.ID,
			"email": mutation.InsertUsersOne.Email,
			"name":  mutation.InsertUsersOne.Name,
		},
	})
	return

}
