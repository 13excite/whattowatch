package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
)

type FilmData struct {
	Title      string
	Genre      string
	PosterLink string
	RatingKp   []string
	RatingImdb []string
	Country    []string
	LinkToKP   string
}

type MinMaxIds struct {
	Min int
	Max int
}

func (ids *MinMaxIds) getMinMaxIds(db *sql.DB) (*MinMaxIds, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	query := "select min(film_id), max(film_id) from films"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	minMaxArr := make([]*MinMaxIds, 1, 1)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&ids.Min, &ids.Max)
		if err != nil {
			return nil, err
		}
		minMaxArr = append(minMaxArr, ids)
	}
	return minMaxArr[1], nil

}

func (f *FilmData) getRandomFilm(db *sql.DB, randomFilmId int) (*FilmData, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	statement := fmt.Sprintf("SELECT title,genre,poster_link,rating_kp,rating_imdb,country,linktokp  FROM films WHERE film_id=%d", randomFilmId)
	rows, err := db.QueryContext(ctx, statement)

	filmsArr := make([]*FilmData, 1, 1)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&f.Title, &f.Genre, &f.PosterLink,
			pq.Array(&f.RatingKp), pq.Array(&f.RatingImdb), pq.Array(&f.Country), &f.LinkToKP)
		if err != nil {
			return nil, err
		}
		filmsArr = append(filmsArr, f)
	}
	return filmsArr[1], nil
}
