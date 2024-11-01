package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"net/http"
)

type oAuthController struct {
	sessionStore sessions.Store
}

type InterfaceOAuthController interface {
	AuthHandler(c *gin.Context)
	CallbackHandler(c *gin.Context)
}

func NewAuthController(sessionStore sessions.Store) InterfaceOAuthController {
	return &oAuthController{
		sessionStore: sessionStore,
	}
}

func (a *oAuthController) AuthHandler(c *gin.Context) {
	providerName := c.Param("provider")

	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider name is required"})
		return
	}

	provider := c.Request.URL.Query()
	provider.Add("provider", "google")
	c.Request.URL.RawQuery = provider.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *oAuthController) CallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in login process"})
		return
	}

	session, err := a.sessionStore.Get(c.Request, "userSession")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in store process"})
	}

	session.Values["user"] = user
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{
		"message":  "authentication successful",
		"user":     user.UserID + ", " + user.Name + ", " + user.Email,
		"userFull": user,
	})
}
