package routes

import (
	"user-api/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetUsersRoutes(r *gin.RouterGroup, c controllers.UserController) {
	r.GET("", c.GetAll())
	r.DELETE("/:id", c.Delete())
	r.GET("/:id", c.GetById())
	r.PUT("/:id", c.Update())
}
