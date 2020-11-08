package main

import (
	"context"
	"database/sql"
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

func (f *FilmData) insertDataToDB(db *sql.DB) error {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	insertQuery := "INSERT INTO films (title, genre, poster_link, rating_kp, rating_imdb, country, linktokp) VALUES($1, $2, $3, $4, $5, $6, $7)"

	_, err = db.ExecContext(ctx, insertQuery, f.Title, f.Genre, f.PosterLink,
		pq.Array(f.RatingKp), pq.Array(f.RatingImdb), pq.Array(f.Country), f.LinkToKP)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ids *MinMaxIds) getMinMaxIds( db *sql.DB) (*MinMaxIds, error) {
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
