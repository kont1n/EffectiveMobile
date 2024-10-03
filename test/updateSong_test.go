package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"EffectiveMobile/internal/models"
)

func TestUpdateSong(t *testing.T) {
	t.Log("Update song")
	// Создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Somewhere I Belong",
		Artist: "Linkin Park",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, song, err := CreateSong(t, testSong, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// Проверяем обновление тестовых данных
	testUpdateSong := models.Song{
		ID:      song.ID,
		Name:    testSong.Name,
		Artist:  testSong.Artist,
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\\nOoh baby, can\nyou hear me moan?\\nYou caught me under false pretenses\\nHow long\nbefore you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set\nmy soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	var detailSong models.Song
	statusCode, detailSong, err = UpdateSong(t, testUpdateSong, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testUpdateSong, detailSong)

	// Удаляем тестовые данные
	statusCode, err = DeleteSong(t, song, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestUpdateSongWithEmptyID(t *testing.T) {
	t.Log("Update song with empty ID")

	var statusCode int
	var err error

	testUpdateSong := models.Song{
		ID:      "",
		Name:    "Supermassive Black Hole",
		Artist:  "Muse",
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\\nOoh baby, can\nyou hear me moan?\\nYou caught me under false pretenses\\nHow long\nbefore you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set\nmy soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, _, err = UpdateSong(t, testUpdateSong, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
