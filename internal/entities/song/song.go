package song

type song struct {
	id       string
	title    string
	duration string
	artist   string
}

type InterfaceSong interface {
	GetId() string
	GetTitle() string
	GetDuration() string
	GetArtist() string
}

func NewSong(id string, title string, duration string, artist string) InterfaceSong {
	return &song{
		id:       id,
		title:    title,
		duration: duration,
		artist:   artist,
	}
}

func (s *song) GetId() string {
	return s.id
}

func (s *song) GetTitle() string {
	return s.title
}

func (s *song) GetDuration() string {
	return s.duration
}

func (s *song) GetArtist() string {
	return s.artist
}
