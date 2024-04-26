package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"app/api/utils"
	db "app/db/sqlc"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	conf *oauth2.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error loading .env file")
	}
	conf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/google/login", s.googleLogin)
	router.GET("/google/callback", s.googleCallback)
	router.GET("/google/logout", s.logout)
}

func (s *Service) googleLogin(c *gin.Context) {
	state := utils.RandToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.Redirect(302, conf.AuthCodeURL(state))
}

func googleUserInfoHandler(token *oauth2.Token, c *gin.Context, s *Service) (*string, error) {
	client := conf.Client(context.Background(), token)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(400, err)
		return nil, err
	}
	if email.StatusCode != 200 {
		return nil, errors.New("invalid access token")
	}
	defer email.Body.Close()

	data, _ := io.ReadAll(email.Body)
	var u User
	json.Unmarshal(data, &u)

	var user_id uuid.UUID
	user, err := s.queries.GetUserByEmail(context.Background(), u.Email)
	if err != nil {
		log.Println(err.Error())
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println(err.Error())
			return nil, err
		}
		createdUser, err := s.queries.CreateUser(context.Background(), db.CreateUserParams{
			Name:    u.Name,
			Email:   u.Email,
			Picture: u.Picture,
		})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		user_id = createdUser
	} else {
		user_id = user.ID
	}
	rand_tok := utils.RandToken()
	sessionErr := s.queries.CreateOrUpdateSession(context.Background(), db.CreateOrUpdateSessionParams{
		User:  user_id,
		Token: rand_tok,
	})
	if sessionErr != nil {
		log.Println(sessionErr.Error())
		return nil, sessionErr
	}
	return &rand_tok, nil
}

func (s *Service) googleCallback(c *gin.Context) {
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.JSON(400, gin.H{"error": "invalid state"})
		return
	}

	tok, err := conf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(400, err)
		return
	}

	rand_tok, err := googleUserInfoHandler(tok, c, s)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(400, err)
		return
	}
	if os.Getenv("GIN_MODE") == "release" {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie("session", *rand_tok, 3600*24*2, "/", os.Getenv("COOKIE_SET_URL"), true, true)
	c.Redirect(302, os.Getenv("LOGIN_REDIRECT_URL"))
}

func (s *Service) logout(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	session := c.MustGet("session").(string)
	err := s.queries.DeleteSession(context.Background(), db.DeleteSessionParams{
		User:  userID,
		Token: session,
	})
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(400, err)
		return
	}
	c.SetCookie("session", "", -1, "/", os.Getenv("COOKIE_SET_URL"), true, true)
	c.Redirect(302, os.Getenv("LOGOUT_REDIRECT_URL"))
}
