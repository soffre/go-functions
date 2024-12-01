package routes

import (
	handlers "hasura/go-functions/Handlers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {

	r.GET("/", handlers.Get_hello)
	r.POST("/create_user", handlers.Create_user)
	r.POST("/name", handlers.Comment_trigger)
}
