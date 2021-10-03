package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"ww/internal/store"
)

func (c *Client) GetRandomFilm(ctx context.Context, randomFilmId int) (*store.FilmData, error) {
	film := new(store.FilmData)

	statement := fmt.Sprintf("SELECT title,genre,poster_link,rating_kp,rating_imdb,country,linktokp  "+
		"FROM films WHERE film_id=%d", randomFilmId)
	rows, err := c.db.QueryContext(ctx, statement)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, wrapError(err)
	}

	filmsArr := make([]*store.FilmData, 1, 1)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&film.Title, &film.Genre, &film.PosterLink,
			pq.Array(&film.RatingKp), pq.Array(&film.RatingImdb), pq.Array(&film.Country), &film.LinkToKP)
		if err != nil {
			return nil, wrapError(err)
		}
		filmsArr = append(filmsArr, film)
	}
	return filmsArr[1], nil
}

func (c *Client) GetMinMaxIds(ctx context.Context) (*store.MinMaxIds, error) {
	minMaxId := new(store.MinMaxIds)

	statement := "select min(film_id), max(film_id) from films"

	err := c.db.GetContext(ctx, minMaxId, statement)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, wrapError(err)
	}
	return minMaxId, nil
}
