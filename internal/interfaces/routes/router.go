package routes

import (
	"github.com/gin-gonic/gin"
	"youtubeapi/internal/interfaces/controllers"
	"youtubeapi/internal/interfaces/middleware/auth"
)

type router struct {
	oauth              controllers.InterfaceOAuthController
	internalAuth       auth.InterfaceAuthenticationMiddleware
	playlistController controllers.InterfacePlaylistController
}

type InterfaceRouter interface {
	InitRoutes(rg *gin.RouterGroup)
}

func NewRouter(oauth controllers.InterfaceOAuthController, internalAuth auth.InterfaceAuthenticationMiddleware, playlistController controllers.InterfacePlaylistController) InterfaceRouter {
	return &router{
		oauth:              oauth,
		internalAuth:       internalAuth,
		playlistController: playlistController,
	}
}

func (r *router) InitRoutes(rg *gin.RouterGroup) {
	rg.GET("/auth/:provider", r.internalAuth.SessionMiddleware, r.oauth.AuthHandler)
	rg.GET("/auth/:provider/callback", r.internalAuth.SessionMiddleware, r.oauth.CallbackHandler)

	authorized := rg.Group("/", r.internalAuth.AuthMiddleware)
	{
		authorized.GET("/playlist", r.playlistController.GetPlaylistHandler)
	}
}
