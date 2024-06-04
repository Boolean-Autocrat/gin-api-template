package middlewares

import (
	"log"
	"net/http"
	"os"
	"time"

	"app/api/utils"
	db "app/db/sqlc"

	"github.com/gin-gonic/gin"
)

func verifySession(c *gin.Context, s *Service) {
	session_token, _ := c.Cookie("session")
	if session_token == "" {
		c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
		return
	}
	session, err := s.queries.GetSession(c, session_token)
	if err != nil {
		log.Print(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	if session.CreatedAt.AddDate(0, 0, 2).Before(time.Now()) {
		new_token := utils.RandToken()
		sessionErr := s.queries.CreateOrUpdateSession(c, db.CreateOrUpdateSessionParams{
			User:  session.User,
			Token: new_token,
		})
		if sessionErr != nil {
			log.Print(sessionErr.Error())
			c.AbortWithStatusJSON(500, gin.H{"message": "internal server error"})
			return
		}
		if os.Getenv("GIN_MODE") == "release" {
			c.SetSameSite(http.SameSiteLaxMode)
		} else {
			c.SetSameSite(http.SameSiteNoneMode)
		}
		c.SetCookie("session", new_token, 3600*24*2, "/", os.Getenv("COOKIE_SET_URL"), true, true)
		session_token = new_token
	}
	c.Set("user_id", session.User)
	c.Set("session", session_token)
	c.Next()
}
