package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log"
	"youtubeapi/configs"
	"youtubeapi/internal/adapters/gateways"
	"youtubeapi/internal/interfaces/controllers"
	"youtubeapi/internal/interfaces/middleware/auth"
	"youtubeapi/internal/interfaces/routes"
	"youtubeapi/internal/useCase"
)

func init() {
	configs.InitEnvConfig()

	goth.UseProviders(
		google.New(configs.EnvConfigs.ClientID, configs.EnvConfigs.SecretKey, "http://localhost:8080/auth/google/callback"),
	)
}

func main() {
	router := gin.Default()

	sessionStore := sessions.NewCookieStore([]byte(configs.EnvConfigs.SessionSecret))
	gothic.Store = sessionStore

	youtubeGateway, err := gateways.NewYoutubeAPIGateway(
		configs.EnvConfigs.ClientID,
		configs.EnvConfigs.SecretKey,
		"http://localhost:8080/auth/google/callback",
	)

	if err != nil {
		log.Fatalf("Error initializing Youtube Gateway: %v", err)
	}

	authMiddleware := auth.NewAuthenticationMiddleware(sessionStore)
	showPlaylistUseCase := useCase.NewShowPlaylistUseCase(youtubeGateway, authMiddleware)
	playlistController := controllers.NewPlaylistController(showPlaylistUseCase)
	oauthController := controllers.NewAuthController(sessionStore)
	routerImpl := routes.NewRouter(oauthController, authMiddleware, playlistController)
	routerImpl.InitRoutes(&router.RouterGroup)
	log.Fatal(router.Run(":8080"))
}
