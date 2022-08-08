package routes

import (
	"user-api/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(r *gin.RouterGroup, c controllers.AuthController) {
	r.POST("/register", c.Register())
	r.POST("/login", c.Login())
}
