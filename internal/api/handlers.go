package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/respond"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/service"
)

var err error

type ApiHandler struct {
	service *service.Service
	loger   *zap.SugaredLogger
	SongAPI
}

type SongAPI interface {
	createSong(writer http.ResponseWriter, request *http.Request)
	readSong(writer http.ResponseWriter, request *http.Request)
	updateSong(writer http.ResponseWriter, request *http.Request)
	deleteSong(writer http.ResponseWriter, request *http.Request)
	getSongInfo(writer http.ResponseWriter, request *http.Request)
	getSongsList(writer http.ResponseWriter, request *http.Request)
	getSongVerses(writer http.ResponseWriter, request *http.Request)
}

func NewHandler(service *service.Service, loger *zap.SugaredLogger) *ApiHandler {
	return &ApiHandler{
		service: service,
		loger:   loger,
	}
}

// createSong : Обработка запроса для создания песни
func (h *ApiHandler) createSong(writer http.ResponseWriter, request *http.Request) {
	var song models.SongRequest
	var result models.SongResponse
	var buf bytes.Buffer

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		h.loger.Errorf("Error reading request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error reading request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		h.loger.Errorf("Error unmarshalling request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error unmarshalling request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if song.Name == "" || song.Artist == "" {
		h.loger.Errorf("Song name, artist not filled")
		h.JSONError(writer, "Song name, artist are required", http.StatusBadRequest, reqID)
		return
	}

	result.ID, err = h.service.CreateSong(song, reqID)
	if err != nil {
		h.loger.Errorf("Error creating song: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error creating song: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusCreated)
}

// readSong : Обработка запроса для получения песни по ее ID
func (h *ApiHandler) readSong(writer http.ResponseWriter, request *http.Request) {
	var result models.Song

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	guid := chi.URLParam(request, "id")
	if guid == "" {
		h.JSONError(writer, "ID is required", http.StatusBadRequest, reqID)
		return
	}

	result, err = h.service.ReadSong(guid, reqID)
	if err != nil {
		h.JSONError(writer, fmt.Sprintf("Error reading song: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// updateSong : Обработка запроса для обновления песни
func (h *ApiHandler) updateSong(writer http.ResponseWriter, request *http.Request) {
	var song models.Song
	var result models.Song
	var buf bytes.Buffer

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		h.loger.Errorf("Error reading request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error reading request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		h.loger.Errorf("Error unmarshalling request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error unmarshalling request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if song.ID == "" || song.Name == "" || song.Artist == "" {
		h.loger.Errorf("Song name, artist not filled")
		h.JSONError(writer, "Song id, name, artist are required", http.StatusBadRequest, reqID)
		return
	}

	result, err = h.service.UpdateSong(song, reqID)
	if err != nil {
		h.loger.Errorf("Error updating song: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error updating song: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// deleteSong : Обработка запроса для удаления песни по ее ID
func (h *ApiHandler) deleteSong(writer http.ResponseWriter, request *http.Request) {
	var result models.SongResponse

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	guid := chi.URLParam(request, "id")
	if guid == "" {
		h.JSONError(writer, "ID is required", http.StatusBadRequest, reqID)
		return
	}

	result, err = h.service.DeleteSong(guid, reqID)
	if err != nil {
		h.JSONError(writer, fmt.Sprintf("Error deleting song: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// getSongInfo : Получение информации о песне
func (h *ApiHandler) getSongInfo(writer http.ResponseWriter, request *http.Request) {
	var song models.SongRequest
	var result models.SongInfoResponse
	var buf bytes.Buffer

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		h.loger.Errorf("Error reading request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error reading request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		h.loger.Errorf("Error unmarshalling request body: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error unmarshalling request body: %v", err.Error()), http.StatusBadRequest, reqID)
		return
	}

	if song.Name == "" || song.Artist == "" {
		h.loger.Errorf("Song name, artist not filled")
		h.JSONError(writer, "Song name, artist are required", http.StatusBadRequest, reqID)
		return
	}

	result, err = h.service.GetSongInfo(song, reqID)
	if err != nil {
		h.loger.Errorf("Error getting song info: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error getting song info: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// getSongsList : Обработка запроса для получения списка песен
func (h *ApiHandler) getSongsList(writer http.ResponseWriter, request *http.Request) {
	var result models.SongsListResponse

	// add pagination and sorting

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	result, err = h.service.GetSongsList(reqID)
	if err != nil {
		h.loger.Errorf("Error getting songs list: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error getting songs list: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// getSongVerses : Обработка запроса для получения куплета песни
func (h *ApiHandler) getSongVerses(writer http.ResponseWriter, request *http.Request) {
	var result models.SongVerseResponse
	var coupletId string

	// add pagination and couplet

	reqID := middleware.GetReqID(request.Context())
	uri := request.RequestURI
	method := request.Method
	h.loger.Debugf("RequestID: %v uri: %v method: %v", reqID, uri, method)

	coupletId = chi.URLParam(request, "сoupletId")
	if coupletId == "" {
		coupletId = defaultСouplet
	}

	result, err = h.service.GetSongVerses(coupletId, reqID)
	if err != nil {
		h.loger.Errorf("Error getting song verses: %v", err)
		h.JSONError(writer, fmt.Sprintf("Error getting song verses: %v", err.Error()), http.StatusInternalServerError, reqID)
		return
	}

	respond.WithJSON(writer, result, http.StatusOK)
}

// JSONError : Обработка ошибок в JSON формате
func (h *ApiHandler) JSONError(w http.ResponseWriter, error string, code int, reqID string) {
	var resp []byte
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	answer := models.ErrorResponse{
		RequestID: reqID,
		Error:     error,
	}
	resp, err = json.Marshal(answer)
	if err != nil {
		h.loger.Errorf("Error marshalling response: %v", err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.loger.Errorf("Error writing response: %v", err)
		return
	}
}
