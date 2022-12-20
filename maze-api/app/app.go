package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/egurnov/maze-api/docs" //nolint:golint
	"github.com/egurnov/maze-api/maze-api/model"
)

type App struct {
	Log *logrus.Logger

	JWTService  model.JWTService
	UserService model.UserService
}

//go:generate swag init -dir ./../../dreamteam --generalInfo ./app/app.go  -o ../../docs

// http://localhost:8080/swagger/index.html

// @title Maze API
// @version 0.1
// @description This API allows to store and solve simple mazes

// @contact.name Alexander Egurnov
// @contact.email alexander.egurnov@gmail.com

// @host localhost:8080
// @BasePath /
// @Schemes http

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization

func (a *App) SetRoutes(r *gin.Engine) {
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, &Message{
			Message: "method not allowed",
		})
	})
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, &Message{
			Message: "not found",
		})
	})
	r.Use(a.renderErrors())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.
		POST("/login", a.Login).
		POST("/user", a.CreateUser)

	// maze := r.Group("/maze", a.AuthorizeJWT())
	// maze.
	// 	GET("", a.GetAllMazes).
	// 	GET(":id", a.GetMaze).
	// 	GET(":id/solution", a.SolveMaze).
	// 	POST("", a.CreateMaze)
}
