package controllers

import (
	"fmt"
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

	fmt.Println(providerName)

	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider name is required"})
		return
	}

	provider := c.Request.URL.Query()
	provider.Add("provider", "google")
	c.Request.URL.RawQuery = provider.Encode()

	gothic.Store = a.sessionStore

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *oAuthController) CallbackHandler(c *gin.Context) {
	gothic.Store = a.sessionStore

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in login process"})
		return
	}

	//providerName := c.Param("provider")

	err = gothic.StoreInSession("token", user.AccessToken, c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in store process"})
		return
	}

	session, err := a.sessionStore.Get(c.Request, "userSession")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in store process"})
	}
	session.Values["userID"] = user.UserID
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{
		"message":  "authentication successful",
		"user":     user.UserID + ", " + user.Name + ", " + user.Email,
		"userFull": user,
	})
}
