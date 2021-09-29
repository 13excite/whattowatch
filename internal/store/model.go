package store

type FilmData struct {
	Title      string
	Genre      string
	PosterLink string `json:"poster_link,omitempty" db:"poster_link"`
	RatingKp   []string `json:"rating_kp,omitempty" db:"rating_kp""`
	RatingImdb []string `json:"rating_imdb,omitempty" db:"rating_imdb"`
	Country    []string
	LinkToKP   string
}

type MinMaxIds struct {
	Min int
	Max int
}
