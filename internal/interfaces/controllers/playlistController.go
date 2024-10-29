package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youtubeapi/internal/useCase"
)

type playlistController struct {
	showPlaylistUseCase useCase.InterfaceShowPlaylistUseCase
}

type InterfacePlaylistController interface {
	GetPlaylistHandler(c *gin.Context)
}

func NewPlaylistController(showPlaylistUseCase useCase.InterfaceShowPlaylistUseCase) InterfacePlaylistController {
	return &playlistController{
		showPlaylistUseCase: showPlaylistUseCase,
	}
}

func (pc *playlistController) GetPlaylistHandler(c *gin.Context) {
	playlistDTO, userDTO, err := pc.showPlaylistUseCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlistDTO": playlistDTO, "userDTO": userDTO})
}
