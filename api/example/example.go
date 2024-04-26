package example

import (
	db "app/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/example", s.getExample)
}

func (s *Service) getExample(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
