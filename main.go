package main

import (
	"log"
	"os"

	"app/api/auth"
	"app/api/example"
	db "app/db/sqlc"
	"app/middlewares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var store = cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	postgres, err := db.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}

	queries := db.New(postgres.DB)

	authService := auth.NewService(queries)
	exampleService := example.NewService(queries)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	if os.Getenv("GIN_MODE") != "release" {
		router.Use(middlewares.CORSMiddleware())
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is healthy!",
		})
	})

	router.Use(sessions.Sessions("goquestsession", store))

	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.NewService(queries).AuthMiddleware())

	exampleGroup := router.Group("")

	authService.RegisterHandlers(authGroup)
	exampleService.RegisterHandlers(exampleGroup)

	router.Run()
}
