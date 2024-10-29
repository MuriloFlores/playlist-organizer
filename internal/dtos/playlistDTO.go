package dtos

type PlaylistDTO struct {
	ID          string `json:"id"`
	PublishedAt string `json:"published_at"`
	ChannelID   string `json:"channel_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
