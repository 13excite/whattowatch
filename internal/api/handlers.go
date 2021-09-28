package api

import (
	"net/http"
	"ww/internal/store"
	"ww/internal/utils"
)

func (s *Server) GetRandomFilm(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	minMaxId, err := s.grStore.GetMinMaxIds(ctx)
	if err != nil {
		errID := RenderErrInternalWithID(w, nil)
		s.logger.Errorw("GetMinMaxIds error", "error", err, "error_id", errID)
		return
	}

	randomNum := utils.RandomRange(minMaxId.Min, minMaxId.Max)

	film, err := s.grStore.GetRandomFilm(ctx, randomNum)
	if err != nil {
		if err == store.ErrNotFound {
			RenderErrResourceNotFound(w, "film")
		} else if serr, ok := err.(*store.Error); ok {
			RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpGet))
		} else {
			errID := RenderErrInternalWithID(w, nil)
			s.logger.Errorw("GetRandomFilm error", "error", err, "error_id", errID)
		}
		return
	}

	RenderJSON(w, http.StatusOK, film)
}
