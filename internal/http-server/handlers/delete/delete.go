package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/wlcmtunknwndth/REST_API/internal/lib/api/response"
	"github.com/wlcmtunknwndth/REST_API/internal/lib/logger/sl"
	"github.com/wlcmtunknwndth/REST_API/internal/storage"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.41.0 --name=URLDelete
type URLDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("didn't get alias")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDelete.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		log.Info("url deleted", slog.String("alias", alias))
	}
}
