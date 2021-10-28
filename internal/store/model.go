package store

type FilmData struct {
	Title      string   `json:"title"`
	Genre      string   `json:"genre"`
	PosterLink string   `json:"poster_link,omitempty" db:"poster_link"`
	RatingKp   []string `json:"rating_kp,omitempty" db:"rating_kp""`
	RatingImdb []string `json:"rating_imdb,omitempty" db:"rating_imdb"`
	Country    []string `json:"country"`
	LinkToKP   string   `json:"link_to_kp"`
}

type MinMaxIds struct {
	Min int
	Max int
}
