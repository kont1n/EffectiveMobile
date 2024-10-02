package service

import (
	"go.uber.org/zap"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/storage"
)

var err error

type Song interface {
}

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
func (s Service) CreateSong(song models.Song, reqID string) (string, error) {
	s.loger.Debugf("RequestID: %v. Creating song in service", reqID)

	song.ID, err = s.store.CreateSong(song, reqID)

	if err != nil {
		s.loger.Errorf("Error creating song: %v", err)
		return "", err
	}

	return song.ID, nil
}
