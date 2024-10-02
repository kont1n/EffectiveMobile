package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/respond"
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

// postSong : Обработка запроса для создания песни
func (h *ApiHandler) postSong(writer http.ResponseWriter, request *http.Request) {
	var song models.SongPostRequest
	var result models.SongPostResponse
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

// getSong : Обработка запроса для получения песни по ее ID
func (h *ApiHandler) getSong(w http.ResponseWriter, r *http.Request) {
	resp := models.SongDetail{
		ID:      "80fbcef2-19d6-49fa-bf54-0ff57a7d0744",
		Name:    "Supermassive Black Hole",
		Artist:  "Muse",
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	respond.WithJSON(w, resp, http.StatusOK)
}

// updateSong : Обработка запроса для обновления песни
func (h *ApiHandler) updateSong(w http.ResponseWriter, r *http.Request) {
	resp := models.SongDetail{
		ID:      "80fbcef2-19d6-49fa-bf54-0ff57a7d0744",
		Name:    "Supermassive Black Hole",
		Artist:  "Muse",
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	respond.WithJSON(w, resp, http.StatusOK)
}

// deleteSong : Обработка запроса для удаления песни по ее ID
func (h *ApiHandler) deleteSong(w http.ResponseWriter, r *http.Request) {
	resp := models.SongDetail{
		ID:      "80fbcef2-19d6-49fa-bf54-0ff57a7d0744",
		Name:    "Supermassive Black Hole",
		Artist:  "Muse",
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	respond.WithJSON(w, resp, http.StatusOK)
}

// getSongsList : Обработка запроса для получения списка песен
func (h *ApiHandler) getSongsList(w http.ResponseWriter, r *http.Request) {
	resp := []models.SongDetail{
		{
			ID:      "80fbcef2-19d6-49fa-bf54-0ff57a7d0744",
			Name:    "Supermassive Black Hole",
			Artist:  "Muse",
			Release: "16.07.2006",
			Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		},
	}
	respond.WithJSON(w, resp, http.StatusOK)
}

// getSongVerses : Обработка запроса для получения куплета песни
func (h *ApiHandler) getSongVerses(w http.ResponseWriter, r *http.Request) {
	resp := models.SongVerseResponse{
		ID:        "80fbcef2-19d6-49fa-bf54-0ff57a7d0744",
		Name:      "Supermassive Black Hole",
		Artist:    "Muse",
		CoupletId: 1,
		Couplet:   "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
	}
	respond.WithJSON(w, resp, http.StatusOK)
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
