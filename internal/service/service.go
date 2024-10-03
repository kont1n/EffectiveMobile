package service

import (
	"strconv"
	"strings"

	"go.uber.org/zap"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/storage"
)

var err error

type Service struct {
	store *storage.Storage
	loger *zap.SugaredLogger
}

func NewService(store *storage.Storage, loger *zap.SugaredLogger) *Service {
	return &Service{
		store: store,
		loger: loger,
	}
}

// CreateSong : Создание песни и вызов сервиса хранилища
func (s Service) CreateSong(song models.SongRequest, reqID string) (string, error) {
	s.loger.Debugf("RequestID: %v. Creating song in service", reqID)
	result := models.SongResponse{}

	result.ID, err = s.store.CreateSong(song, reqID)
	if err != nil {
		s.loger.Errorf("Error creating song: %v", err)
		return "", err
	}

	return result.ID, nil
}

// ReadSong : Получение песни по ее ID и вызов сервиса хранилища
func (s Service) ReadSong(guid string, reqID string) (models.Song, error) {
	s.loger.Debugf("RequestID: %v. Creating song in service", reqID)
	result := models.Song{}

	result, err = s.store.ReadSong(guid, reqID)
	if err != nil {
		s.loger.Errorf("Error getting song: %v", err)
		return result, err
	}

	return result, nil
}

// UpdateSong : Обновление песни и вызов сервиса хранилища
func (s Service) UpdateSong(song models.Song, reqID string) (models.Song, error) {
	s.loger.Debugf("RequestID: %v. Updating song in service", reqID)
	result := models.Song{}

	result, err = s.store.UpdateSong(song, reqID)
	if err != nil {
		s.loger.Errorf("Error updating song: %v", err)
		return result, err
	}

	return result, nil
}

// DeleteSong : Удаление песни и вызов сервиса хранилища
func (s Service) DeleteSong(guid string, reqID string) (models.SongResponse, error) {
	s.loger.Debugf("RequestID: %v. Deleting song in service", reqID)
	result := models.SongResponse{}

	result, err = s.store.DeleteSong(guid, reqID)
	if err != nil {
		s.loger.Errorf("Error deleting song: %v", err)
		return result, err
	}

	return result, nil
}

// GetSongInfo : Получение информации о песне и вызов сервиса хранилища
func (s Service) GetSongInfo(song models.SongRequest, reqID string) (models.SongInfoResponse, error) {
	s.loger.Debugf("RequestID: %v. Getting song info in service", reqID)
	result := models.SongInfoResponse{}

	result, err = s.store.GetSongInfo(song, reqID)
	if err != nil {
		s.loger.Errorf("Error getting song info: %v", err)
		return result, err
	}

	return result, nil
}

// GetSongsList : Получение списка песен и вызов сервиса хранилища
func (s Service) GetSongsList(reqID string) (models.SongsListResponse, error) {
	s.loger.Debugf("RequestID: %v. Getting songs list in service", reqID)
	result := models.SongsListResponse{}

	result, err = s.store.GetSongsList(reqID)
	if err != nil {
		s.loger.Errorf("Error getting songs list: %v", err)
		return result, err
	}

	return result, nil
}

// GetSongVerses : Получение куплета песни и вызов сервиса хранилища
func (s Service) GetSongVerses(coupletId string, reqID string) (models.SongVerseResponse, error) {
	var couplet int
	var couplets []string
	s.loger.Debugf("RequestID: %v. Getting song verses in service", reqID)
	result := models.SongVerseResponse{}

	couplet, err = strconv.Atoi(coupletId)
	if err != nil {
		s.loger.Errorf("Error converting coupletId to int: %v", err)
		return result, err
	}

	result, err = s.store.GetSongVerses(reqID)
	if err != nil {
		s.loger.Errorf("Error getting song verses: %v", err)
		return result, err
	}

	couplets = strings.Split(result.Couplet, "\n\n")
	if len(couplets) < couplet {
		s.loger.Errorf("Error getting song verses: coupletId is out of range")
		return result, err
	}

	result.Couplet = couplets[couplet-1]

	return result, nil
}
