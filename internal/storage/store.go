package storage

import (
	"context"
	"database/sql"
	"time"

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
	SongStorage
}
type SongStorage interface {
	CreateSong(song models.SongRequest, reqID string) (string, error)
	ReadSong(guid string, reqID string) (models.Song, error)
	UpdateSong(song models.Song, reqID string) (models.Song, error)
	DeleteSong(guid string, reqID string) (models.SongResponse, error)
	GetSongInfo(song models.SongRequest, reqID string) (models.SongInfoResponse, error)
	GetSongsList(reqID string) (models.SongsListResponse, error)
	GetSongCouplet(guid string, id string) (string, error)
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
func (s Storage) ReadSong(guid string, reqID string) (result models.Song, err error) {
	s.loger.Debugf("RequestID: %v. Reading song from the database", reqID)
	resultDB := models.SongDB{}

	err = s.db.QueryRow(context.Background(), "SELECT id, song_name, artist_name, release_date, song_text, link FROM public.songs WHERE id = $1",
		guid).Scan(&resultDB.ID, &resultDB.Name, &resultDB.Artist, &resultDB.Release, &resultDB.Text, &resultDB.Link)

	if err != nil {
		s.loger.Errorf("Error getting song from the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("SongDB: %v", resultDB)

	var releaseDate, text, link string

	if resultDB.Release.Valid {
		releaseDate = resultDB.Release.Time.Format("02.01.2006")
	} else {
		releaseDate = ""
	}
	if resultDB.Text.Valid {
		text = resultDB.Text.String
	} else {
		text = ""
	}
	if resultDB.Link.Valid {
		link = resultDB.Link.String
	} else {
		link = ""
	}

	result = models.Song{
		ID:      resultDB.ID,
		Name:    resultDB.Name,
		Artist:  resultDB.Artist,
		Release: releaseDate,
		Text:    text,
		Link:    link,
	}

	s.loger.Debugf("Song: %v", result)

	s.loger.Debugf("RequestID: %v. Song read from the database", reqID)
	return result, nil
}

// UpdateSong : Обновление песни в базе данных
func (s Storage) UpdateSong(song models.Song, reqID string) (result models.Song, err error) {
	s.loger.Debugf("RequestID: %v. Updating song in the database", reqID)
	resultDB := models.SongDB{}
	var releaseDate time.Time
	var timeVaild, textValid, linkValid bool

	if song.Release != "" {
		releaseDate, _ = time.Parse("02.01.2006", song.Release)
		timeVaild = true
	} else {
		timeVaild = false
		releaseDate = time.Time{}
	}

	if song.Text != "" {
		textValid = true
	} else {
		textValid = false
	}

	if song.Link != "" {
		linkValid = true
	} else {
		linkValid = false
	}

	songDB := models.SongDB{
		ID:      song.ID,
		Name:    song.Name,
		Artist:  song.Artist,
		Release: sql.NullTime{Time: releaseDate, Valid: timeVaild},
		Text:    sql.NullString{String: song.Text, Valid: textValid},
		Link:    sql.NullString{String: song.Link, Valid: linkValid},
	}

	err = s.db.QueryRow(context.Background(), "UPDATE public.songs SET song_name= $1, artist_name= $2, release_date= $3, song_text= $4, link= $5  WHERE id= $6  RETURNING id,song_name,artist_name,release_date,song_text,link",
		songDB.Name, songDB.Artist, songDB.Release, songDB.Text, songDB.Link, songDB.ID).Scan(&resultDB.ID, &resultDB.Name, &resultDB.Artist, &resultDB.Release, &resultDB.Text, &resultDB.Link)

	var rlsDate, text, link string

	if resultDB.Release.Valid {
		rlsDate = resultDB.Release.Time.Format("02.01.2006")
	} else {
		rlsDate = ""
	}
	if resultDB.Text.Valid {
		text = resultDB.Text.String
	} else {
		text = ""
	}
	if resultDB.Link.Valid {
		link = resultDB.Link.String
	} else {
		link = ""
	}

	result = models.Song{
		ID:      resultDB.ID,
		Name:    resultDB.Name,
		Artist:  resultDB.Artist,
		Release: rlsDate,
		Text:    text,
		Link:    link,
	}

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

	err = s.db.QueryRow(context.Background(), "DELETE FROM public.songs WHERE id = $1::uuid RETURNING id",
		guid).Scan(&result.ID)

	if err != nil {
		s.loger.Errorf("Error deleting song in the database: %v", err.Error())
		return result, err
	}

	s.loger.Debugf("RequestID: %v. Song deleted in the database", reqID)
	return result, nil
}

// GetSongInfo : Получение информации о песни в базе данных
func (s Storage) GetSongInfo(song models.SongRequest, reqID string) (result models.SongInfoResponse, err error) {
	s.loger.Debugf("RequestID: %v. Reading song info from the database", reqID)
	resultDB := models.SongInfoResponseDB{}

	err = s.db.QueryRow(context.Background(), "SELECT release_date, song_text, link FROM public.songs WHERE song_name = $1 AND artist_name = $2",
		song.Name, song.Artist).Scan(&resultDB.Release, &resultDB.Text, &resultDB.Link)

	if err != nil {
		s.loger.Errorf("Error getting song info from the database: %v", err.Error())
		return result, err
	}

	var rlsDate, text, link string

	if resultDB.Release.Valid {
		rlsDate = resultDB.Release.Time.Format("02.01.2006")
	} else {
		rlsDate = ""
	}
	if resultDB.Text.Valid {
		text = resultDB.Text.String
	} else {
		text = ""
	}
	if resultDB.Link.Valid {
		link = resultDB.Link.String
	} else {
		link = ""
	}

	result = models.SongInfoResponse{
		Release: rlsDate,
		Text:    text,
		Link:    link,
	}

	s.loger.Info("Song info: ", result)

	s.loger.Debugf("RequestID: %v. Song info read from the database", reqID)
	return result, nil
}

// GetSongsList : Получение списка песен в базе данных
func (s Storage) GetSongsList(reqID string) (models.SongsListResponse, error) {
	var rows pgx.Rows

	s.loger.Debugf("RequestID: %v. Reading songs list from the database", reqID)
	result := models.SongsListResponse{}

	limit := 10

	rows, err = s.db.Query(context.Background(), "SELECT id, song_name, artist_name, release_date, song_text, link FROM public.songs ORDER BY id DESC LIMIT $1", limit)
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

// GetSongCouplet : Получение куплета песни в базе данных
func (s Storage) GetSongCouplet(guid string, id string) (result string, err error) {
	s.loger.Debugf("RequestID: %v. Reading song verses from the database with ID=", id)

	var resultDB sql.NullString

	err = s.db.QueryRow(context.Background(), "SELECT song_text FROM public.songs WHERE id = $1::uuid", guid).Scan(&resultDB)

	if err != nil {
		s.loger.Errorf("Error getting song verses from the database: %v", err.Error())
		return result, err
	}

	if resultDB.Valid {
		result = resultDB.String
	} else {
		result = ""
	}

	s.loger.Info("Text: ", result)

	s.loger.Debugf("RequestID: %v. Song verses read from the database", id)
	return result, nil
}
