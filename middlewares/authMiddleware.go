package middlewares

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

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/auth/google/login" || c.Request.URL.Path == "/auth/google/callback" {
			c.Next()
			return
		}
		verifySession(c, s)
		c.Next()
	}
}
