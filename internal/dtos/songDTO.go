package dtos

type SongDTO struct {
	ID          string `json:"id"`
	PublishedAt string `json:"publishedAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoID     string `json:"videoId"`
}
