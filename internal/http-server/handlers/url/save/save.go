package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shorteener/internal/lib/api/response"
	"url-shorteener/internal/lib/logger/sl"
	"url-shorteener/internal/lib/random"
	"url-shorteener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	aliasLength = 4
)

type Request struct {
	URL		string `json:"url" validate:"required,url"`
	Alias	string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias	string `json:"alias,omitempty"`
}

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
			log.Error("failed to decode reques body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		
		if err:=validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			
			log.Error("invalid req", sl.Err(err))

			render.JSON(w, r, response.Error("invalid req"))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		_, err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to save"))
			return
		}

		render.JSON(w, r, Response{
			Response: 	response.OK(),
			Alias:		alias,
		})
	}
}