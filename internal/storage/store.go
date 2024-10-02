package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"EffectiveMobile/internal/models"
)

var (
	err error
)

type Storage struct {
	db    *pgxpool.Pool
	loger *zap.SugaredLogger
}

func NewStorage(db *pgxpool.Pool, loger *zap.SugaredLogger) *Storage {
	return &Storage{
		db:    db,
		loger: loger,
	}
}

// CreateSong : Создание песни в базе данных
func (s Storage) CreateSong(song models.Song, reqID string) (string, error) {
	s.loger.Debugf("RequestID: %v. Creating song in the database", reqID)

	err = s.db.QueryRow(context.Background(), "INSERT INTO public.songs (song_name, artist_name) VALUES ($1, $2) RETURNING id",
		song.Name, song.Artist).Scan(&song.ID)

	if err != nil {
		s.loger.Errorf("Error creating song in the database: %v", err.Error())
		return "", err
	}

	s.loger.Debugf("RequestID: %v. Song created in the database", reqID)
	return song.ID, nil
}
