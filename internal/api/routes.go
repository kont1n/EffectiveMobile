package api

import (
	"log"
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

	router.Get("/topic/{id}", func(w http.ResponseWriter, r *http.Request) {
		resp := models.Topic{
			Namespace: "example",
			Topic:     "topic",
			Private:   false,
			ViewCount: 412,
		}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	router.Get("/topics", func(w http.ResponseWriter, r *http.Request) {
		resp := models.TopicsGetResponse{
			Topics: []models.TopicRecord{
				{
					ID: "testId",
					Topic: models.Topic{
						Namespace: "example",
						Topic:     "topic",
						Private:   false,
						ViewCount: 412,
					},
				},
			},
		}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	router.Post("/topics", func(w http.ResponseWriter, r *http.Request) {
		resp := models.TopicsPostResponse{ID: "123"}
		respond.WithJSON(w, resp, http.StatusOK)
	})

	// Create the API definition.
	api := rest.NewAPI("Music Store API")

	// Create the routes and parameters of the Router in the REST API definition with an
	// adapter, or do it manually.
	chiadapter.Merge(api, router)

	// It's possible to customise the OpenAPI schema for each type.
	api.RegisterModel(rest.ModelOf[respond.Error](), rest.WithDescription("Standard JSON error"), func(s *openapi3.Schema) {
		status := s.Properties["statusCode"]
		status.Value.WithMin(100).WithMax(600)
	})

	// Document the routes.
	api.Get("/topic/{id}").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.TopicsGetResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[respond.Error]())

	api.Get("/topics").
		HasResponseModel(http.StatusOK, rest.ModelOf[models.TopicsGetResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[respond.Error]())

	api.Post("/topics").
		HasRequestModel(rest.ModelOf[models.TopicsPostRequest]()).
		HasResponseModel(http.StatusOK, rest.ModelOf[models.TopicsPostResponse]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[respond.Error]())

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
		log.Fatalf("failed to create swagger UI handler: %v", err)
	}
	router.Handle("/swagger-ui*", ui)

	router.Post("/api/song", h.postSong)

	return router
}
