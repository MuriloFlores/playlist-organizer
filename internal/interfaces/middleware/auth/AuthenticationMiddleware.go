package auth

import (
	"errors"
	"fmt"
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
	GetSessionData(c *gin.Context) (map[interface{}]interface{}, error)
	AuthMiddleware(c *gin.Context)
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
	/*

		store, ok := c.Get("sessionStore")
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"Error": "Session Store Not Found"})
				c.Abort()
				return
			}

			sessionStore, ok := store.(sessions.CookieStore)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Session Store"})
				c.Abort()
				return
			}

			session, err := sessionStore.Get(c.Request, "userSession")
			if err != nil || (session.Values["userID"] == nil && session.Values["userID"] == "") {
				c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
				c.Abort()
				return
			}

	*/

	c.Get()

	fmt.Println(sesao)

	if err := s.extendAge(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "refresh store error"})
		c.Abort()
		return
	}

	c.Next()
}

func (s *authenticationMiddleware) extendAge(c *gin.Context) error {
	store, ok := c.Get("sessionStore")
	if !ok {
		return errors.New("sessionStore not found in context")
	}

	sessionStore, ok := store.(sessions.CookieStore)
	if !ok {
		return errors.New("invalid session store type")
	}

	sessionStore.Options = &sessions.Options{
		MaxAge:   172800,
		HttpOnly: true,
	}

	return nil
}

func (s *authenticationMiddleware) GetSessionData(c *gin.Context) (map[interface{}]interface{}, error) {
	store, ok := c.Get("sessionStore")
	if !ok {
		return nil, errors.New("sessionStore not found in context")
	}

	sessionStore, ok := store.(sessions.Store)
	if !ok {
		return nil, errors.New("invalid session store type")
	}

	session, err := sessionStore.Get(c.Request, "userSession")
	if err != nil {
		return nil, err
	}

	return session.Values, nil
}
