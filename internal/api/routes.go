package api

import (
	"net/http"

	"github.com/a-h/respond"
	"github.com/a-h/rest"
	"github.com/a-h/rest/chiadapter"
	"github.com/a-h/rest/swaggerui"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"EffectiveMobile/internal/models"
)

var (
	spec *openapi3.T
)

func (h *ApiHandler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(h.LogAPI)

	router.Post("/api/song", h.createSong)
	router.Put("/api/song", h.updateSong)
	router.Delete("/api/song/{id}", h.deleteSong)
	router.Get("/api/song/info", h.getSongInfo)
	router.Get("/api/song/{id}", h.readSong)
	router.Get("/api/song/{id}/couplet", h.getSongCouplet)

	router.Route("/api/songs", func(r chi.Router) {
		r.Use(h.Sorting)
		r.Use(h.Filtering)
		r.Use(h.Pagination)
		r.Get("/", h.getSongsList)
	})

	// Create the API definition.
	api := rest.NewAPI("Music Store API")

	// Create the routes and parameters of the Router in the REST API definition with an
	// adapter, or do it manually.
	err = chiadapter.Merge(api, router)
	if err != nil {
		h.loger.Errorf("Failed to merge router to api: %v", err)
	}

	// It's possible to customise the OpenAPI schema for each type.
	api.RegisterModel(rest.ModelOf[respond.Error](), rest.WithDescription("Standard JSON error"), func(s *openapi3.Schema) {
		status := s.Properties["statusCode"]
		status.Value.WithMin(100).WithMax(600)
	})

	// Document the routes.
	api.Post("/api/song").
		HasRequestModel(rest.ModelOf[models.SongRequest]()).
		HasResponseModel(http.StatusCreated, rest.ModelOf[models.SongResponse]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[models.ErrorResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Get("/api/song/{id}").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.Song]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Put("/api/song").
		HasRequestModel(rest.ModelOf[models.Song]()).
		HasResponseModel(http.StatusOK, rest.ModelOf[models.Song]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[models.ErrorResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Delete("/api/song/{id}").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.SongResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Get("/api/song/info").
		HasRequestModel(rest.ModelOf[models.SongRequest]()).
		HasResponseModel(http.StatusOK, rest.ModelOf[models.SongInfoResponse]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[models.ErrorResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Get("/api/song/{id}/couplet").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.SongVerseResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	api.Get("/api/songs/").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.SongsListResponse]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[models.ErrorResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[models.ErrorResponse]())

	// Create the spec.
	spec, err = api.Spec()
	if err != nil {
		h.loger.Errorf("failed to create spec: %v", err)
	}

	spec.Info.Version = "v1.0.0"
	spec.Info.Description = "Описание интеграционных сервисов для работы с хранилищем песен"

	// Attach the UI handler.
	var ui http.Handler
	ui, err = swaggerui.New(spec)
	if err != nil {
		h.loger.Errorf("failed to create swagger UI handler: %v", err)
	}
	router.Handle("/swagger-ui*", ui)

	return router
}
