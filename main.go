package main

import (
	"context"
	"user-api/controllers/v1"
	database "user-api/databases"
	docs "user-api/docs"
	"user-api/repositories"
	routes "user-api/routes"
	service "user-api/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func main() {
	//init mongo connection
	ctx := context.TODO()
	mongoClient := database.MongoInit(&ctx)
	defer mongoClient.Disconnect(ctx)
	userDb := mongoClient.Database("user-api")

	//init repositories
	userRepo := repositories.NewUserMongo(userDb.Collection("users"), ctx)

	//init userService
	userSvc := service.NewUser(userRepo)

	//init controller
	userController := controllers.NewUserJson(userSvc)
	authController := controllers.NewAuth(userSvc)

	//init v1 router
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"
	v1 := router.Group("/v1")

	//set routes
	userGroup := v1.Group("/users")
	userGroup.Use(authController.VerifyToken())
	routes.SetUsersRoutes(userGroup, userController)
	routes.SetAuthRoutes(v1.Group("/auth"), authController)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8082")
}
