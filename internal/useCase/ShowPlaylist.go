package useCase

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
	"youtubeapi/internal/adapters/gateways"
	"youtubeapi/internal/dtos"
	"youtubeapi/internal/interfaces/middleware/auth"
)

type showPlaylistUseCase struct {
	youtubeGateway gateways.InterfaceYoutubeAPIGateway
	authMiddleware auth.InterfaceAuthenticationMiddleware
}

type InterfaceShowPlaylistUseCase interface {
	Execute(c *gin.Context) (*dtos.PlaylistDTO, *dtos.UserDTO, error)
}

func NewShowPlaylistUseCase(youtubeGateway gateways.InterfaceYoutubeAPIGateway, authMiddleware auth.InterfaceAuthenticationMiddleware) InterfaceShowPlaylistUseCase {
	return &showPlaylistUseCase{
		youtubeGateway: youtubeGateway,
		authMiddleware: authMiddleware,
	}
}

func (uc *showPlaylistUseCase) Execute(c *gin.Context) (*dtos.PlaylistDTO, *dtos.UserDTO, error) {
	userData, err := uc.authMiddleware.GetUserSessionData(c)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting session data: %v", err)
	}

	user, ok := userData["user"].(goth.User)
	if !ok {
		return nil, nil, fmt.Errorf("error access token not found: %v", err)
	}

	oauthToken := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.ExpiresAt,
		ExpiresIn:    172800,
	}

	playlist, err := uc.youtubeGateway.GetPlaylists(c.Request.Context(), oauthToken)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting playlists: %v", oauthToken)
	}

	userDTO, err := uc.youtubeGateway.GetUserInfo(c.Request.Context(), oauthToken)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user info: %v", err)
	}

	playlistDTO := &dtos.PlaylistDTO{
		ID:          playlist.GetId(),
		PublishedAt: playlist.GetPublishedAt(),
		ChannelID:   playlist.GetChannelId(),
		Title:       playlist.GetTitle(),
		Description: playlist.GetDescription(),
	}

	return playlistDTO, userDTO, nil
}
