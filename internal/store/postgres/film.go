package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ww/internal/store"
)

func (c *Client) GetRandomFilm(ctx context.Context, randomFilmId int) (*store.FilmData, error) {
	film := new(store.FilmData)
	statement := fmt.Sprintf("SELECT title,genre,poster_link,rating_kp,rating_imdb,country,linktokp "+
		"FROM films WHERE film_id=%d", randomFilmId)
	err := c.db.GetContext(ctx, film, statement)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, wrapError(err)
	}
	return film, nil
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
