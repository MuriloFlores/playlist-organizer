package playlist

import (
	"youtubeapi/internal/entities/song"
)

type playlist struct {
	id          string
	channelId   string
	title       string
	description string
	publishedAt string
	songs       []song.InterfaceSong
}

type InterfacePlaylist interface {
	GetId() string
	GetChannelId() string
	GetTitle() string
	GetDescription() string
	GetPublishedAt() string
	GetSongs() []song.InterfaceSong
}

func NewPlaylist(id string, channelId string, title string, description string, publishedAt string, songs []song.InterfaceSong) InterfacePlaylist {
	return &playlist{
		id:          id,
		channelId:   channelId,
		title:       title,
		description: description,
		publishedAt: publishedAt,
		songs:       songs,
	}
}

func (p *playlist) GetId() string {
	return p.id
}

func (p *playlist) GetChannelId() string {
	return p.channelId
}

func (p *playlist) GetTitle() string {
	return p.title
}

func (p *playlist) GetDescription() string {
	return p.description
}

func (p *playlist) GetPublishedAt() string {
	return p.publishedAt
}

func (p *playlist) GetSongs() []song.InterfaceSong {
	return p.songs
}
