package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
	"youtubeapi/configs"
)

type authenticationMiddleware struct {
	store sessions.Store
}

type InterfaceAuthenticationMiddleware interface {
	SessionMiddleware(c *gin.Context)
	AuthMiddleware(c *gin.Context)
	GetUserSessionData(c *gin.Context) (map[interface{}]interface{}, error)
}

func NewAuthenticationMiddleware(store sessions.Store) InterfaceAuthenticationMiddleware {
	return &authenticationMiddleware{
		store: store,
	}
}

func (s *authenticationMiddleware) SessionMiddleware(c *gin.Context) {
	if s.store == nil {
		s.store = sessions.NewCookieStore([]byte(configs.EnvConfigs.SessionSecret))
	}

	c.Set("sessionStore", s.store)
	c.Next()
}

func (s *authenticationMiddleware) AuthMiddleware(c *gin.Context) {
	store, err := s.store.Get(c.Request, "userSession")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Session Store Not Found"})
		c.Abort()
		return
	}

	if store.Values["AccessToken"] == nil && store.Values["AccessToken"] == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}

func (s *authenticationMiddleware) extendAge(c *gin.Context) error {
	store, err := s.store.Get(c.Request, "userSession")
	if err != nil {
		return errors.New("sessionStore not found in context")
	}

	store.Options = &sessions.Options{
		MaxAge:   172800,
		HttpOnly: true,
	}

	return nil
}

func (s *authenticationMiddleware) GetUserSessionData(c *gin.Context) (map[interface{}]interface{}, error) {
	store, err := s.store.Get(c.Request, "userSession")
	if err != nil {
		return nil, err
	}

	return store.Values, nil
}
