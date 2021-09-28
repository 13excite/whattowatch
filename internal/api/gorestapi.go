package api

import (
	"context"
	"ww/internal/store"
)

type GRStore interface {
	GetRandomFilm(ctx context.Context, randomFilmId int) (*store.FilmData, error)
	GetMinMaxIds(ctx context.Context) (*store.MinMaxIds, error)
}
