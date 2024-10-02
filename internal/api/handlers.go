package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/service"
)

var err error

type ApiHandler struct {
	service *service.Service
	loger   *zap.SugaredLogger
}

func NewHandler(service *service.Service, loger *zap.SugaredLogger) *ApiHandler {
	return &ApiHandler{
		service: service,
		loger:   loger,
	}
}

func (h *ApiHandler) postSong(writer http.ResponseWriter, request *http.Request) {
	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	var song models.Song
	var buf bytes.Buffer

	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		h.loger.Errorf("Error reading request body: %v", err)
		http.Error(writer, "error reading request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		h.loger.Errorf("Error unmarshalling request body: %v", err)
		http.Error(writer, "error unmarshalling request body", http.StatusBadRequest)
		return
	}

	if song.Name == "" || song.Artist == "" {
		h.loger.Errorf("Song name, artist not filled")
		http.Error(writer, "song name, artist are required", http.StatusBadRequest)
		return
	}

	song.ID, err = h.service.CreateSong(song, reqID)
	if err != nil {
		h.loger.Errorf("Error creating song: %v", err)
		http.Error(writer, "error creating song", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)

	var resp []byte

	resp, err = json.Marshal(song.ID)
	if err != nil {
		h.loger.Errorf("Error marshalling response: %v", err)
		http.Error(writer, "error marshalling response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(resp)
	if err != nil {
		h.loger.Errorf("Error writing response: %v", err)
		http.Error(writer, "error writing response", http.StatusInternalServerError)
		return
	}
}
