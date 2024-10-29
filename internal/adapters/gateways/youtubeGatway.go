package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"youtubeapi/internal/dtos"
	"youtubeapi/internal/entities/playlist"
	"youtubeapi/internal/entities/song"
)

type youtubeAPIGateway struct {
	oauth2Config   *oauth2.Config
	youtubeService *youtube.Service
}

type InterfaceYoutubeAPIGateway interface {
	GetYoutubeService(ctx context.Context, token *oauth2.Token) (*youtube.Service, error)
	GetPlaylists(ctx context.Context, token *oauth2.Token) (playlist.InterfacePlaylist, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*dtos.UserDTO, error)
}

func NewYoutubeAPIGateway(clientID, clientSecret, redirectURL string) (InterfaceYoutubeAPIGateway, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/youtube",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &youtubeAPIGateway{oauth2Config: config}, nil
}

func (y *youtubeAPIGateway) GetYoutubeService(ctx context.Context, token *oauth2.Token) (*youtube.Service, error) {
	if y.youtubeService == nil {
		service, err := youtube.NewService(ctx, option.WithTokenSource(y.oauth2Config.TokenSource(ctx, token)))
		if err != nil {
			return nil, fmt.Errorf("create youtube service error: %w", err)
		}
		y.youtubeService = service
	}

	return y.youtubeService, nil
}

func (y *youtubeAPIGateway) GetPlaylists(ctx context.Context, token *oauth2.Token) (playlist.InterfacePlaylist, error) {
	service, err := y.GetYoutubeService(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("get youtube service error: %w", err)
	}

	call := service.Playlists.List([]string{"snippet,contentDetails"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("search playlists error: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("no youtube playlist found")
	}

	RespPlaylist := response.Items[0]
	return playlist.NewPlaylist(
		RespPlaylist.Id,
		RespPlaylist.Snippet.ChannelId,
		RespPlaylist.Snippet.Title,
		RespPlaylist.Snippet.Description,
		RespPlaylist.Snippet.PublishedAt,
		[]song.InterfaceSong{song.NewSong("u", "u", "u", "u")},
	), nil
}

func (y *youtubeAPIGateway) GetUserInfo(ctx context.Context, token *oauth2.Token) (*dtos.UserDTO, error) {
	client := y.oauth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("get user info error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info error: %s", resp.Status)
	}

	var userInfo dtos.UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info %w: ", err)
	}

	return &userInfo, nil
}
