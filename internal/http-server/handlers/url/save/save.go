package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "github.com/wlcmtunknwndth/REST_API/internal/lib/api/response"
	"github.com/wlcmtunknwndth/REST_API/internal/lib/logger/sl"
	"github.com/wlcmtunknwndth/REST_API/internal/lib/random"
	"github.com/wlcmtunknwndth/REST_API/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// omitempty -- если пустая строка, на в json не указывается

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to cfg
const aliasLength = 4

//go:generate go run github.com/vektra/mockery/v2@v2.41.0 --name=URLSaver

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			//render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		// TODO: fix collisions for alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		//render.JSON(w, r, Response{
		//	Response: resp.OK(),
		//	Alias:    alias,
		//})
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
