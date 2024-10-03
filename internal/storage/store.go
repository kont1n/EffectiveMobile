package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
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
func (s Storage) CreateSong(song models.SongRequest, reqID string) (string, error) {
	s.loger.Debugf("RequestID: %v. Creating song in the database", reqID)
	result := models.SongResponse{}

	err = s.db.QueryRow(context.Background(), "INSERT INTO public.songs (song_name, artist_name) VALUES ($1, $2) RETURNING id",
		song.Name, song.Artist).Scan(&result.ID)

	if err != nil {
		s.loger.Errorf("Error creating song in the database: %v", err.Error())
		return "", err
	}

	s.loger.Debugf("RequestID: %v. Song created in the database", reqID)
	return result.ID, nil
}

// ReadSong : Получение песни по ее ID из базы данных
func (s Storage) ReadSong(guid string, reqID string) (models.Song, error) {
	s.loger.Debugf("RequestID: %v. Reading song from the database", reqID)
	result := models.Song{}

	err = s.db.QueryRow(context.Background(), "SELECT id, song_name, artist_name, release_date, text, link FROM public.songs WHERE id = $1",
		guid).Scan(&result.ID, &result.Name, &result.Artist, &result.Release, &result.Text, &result.Link)

	if err != nil {
		s.loger.Errorf("Error getting song from the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song read from the database", reqID)
	return result, nil
}

// UpdateSong : Обновление песни в базе данных
func (s Storage) UpdateSong(song models.Song, reqID string) (models.Song, error) {
	s.loger.Debugf("RequestID: %v. Updating song in the database", reqID)
	result := models.Song{}

	err = s.db.QueryRow(context.Background(), "UPDATE song_name, artist_name, release_date, text, link FROM public.songs WHERE id = $1",
		song.ID).Scan(&result.ID, &result.Name, &result.Artist, &result.Release, &result.Text, &result.Link)

	if err != nil {
		s.loger.Errorf("Error updating song in the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song updated in the database", reqID)
	return result, nil
}

// DeleteSong : Удаление песни в базе данных
func (s Storage) DeleteSong(guid string, reqID string) (models.SongResponse, error) {
	s.loger.Debugf("RequestID: %v. Deleting song in the database", reqID)
	result := models.SongResponse{}

	err = s.db.QueryRow(context.Background(), "DELETE FROM public.songs WHERE id = $1",
		guid).Scan(&result.ID)

	if err != nil {
		s.loger.Errorf("Error deleting song in the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song deleted in the database", reqID)
	return result, nil
}

// GetSongInfo : Получение информации о песни в базе данных
func (s Storage) GetSongInfo(song models.SongRequest, reqID string) (models.SongInfoResponse, error) {
	s.loger.Debugf("RequestID: %v. Reading song info from the database", reqID)
	result := models.SongInfoResponse{}

	err = s.db.QueryRow(context.Background(), "SELECT release_date, text, link FROM public.songs WHERE song_name = $1 AND artist_name = $2",
		song.Name, song.Artist).Scan(&result.Release, &result.Text, &result.Link)

	if err != nil {
		s.loger.Errorf("Error getting song info from the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song info read from the database", reqID)
	return result, nil
}

// GetSongsList : Получение списка песен в базе данных
func (s Storage) GetSongsList(reqID string) (models.SongsListResponse, error) {
	var rows pgx.Rows

	s.loger.Debugf("RequestID: %v. Reading songs list from the database", reqID)
	result := models.SongsListResponse{}

	limit := 10

	rows, err = s.db.Query(context.Background(), "SELECT id, song_name, artist_name, release_date, text, link FROM public.songs ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		s.loger.Errorf("Error getting songs list from the database: %v", err.Error())
		return result, err
	}

	result.Songs, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.Song])
	if err != nil {
		s.loger.Errorf("Error collecting rows: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Songs list read from the database", reqID)
	return result, nil
}

// GetSongVerses : Получение куплета песни в базе данных
func (s Storage) GetSongVerses(id string) (models.SongVerseResponse, error) {
	s.loger.Debugf("RequestID: %v. Reading song verses from the database", id)
	result := models.SongVerseResponse{}

	err = s.db.QueryRow(context.Background(), "SELECT text FROM public.song_verses WHERE song_id = $1",
		id).Scan(&result.Couplet)

	if err != nil {
		s.loger.Errorf("Error getting song verses from the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song verses read from the database", id)
	return result, nil
}
