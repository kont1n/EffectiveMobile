package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"EffectiveMobile/internal/models"
)

func TestGetSongInfo(t *testing.T) {
	t.Log("Get song info")
	// Создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Raging Bull",
		Artist: "Paul Velchev",
	}

	statusCode, song, err := CreateSong(t, testSong)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// Обновляем тестовые данные
	testUpdateSong := models.Song{
		ID:      song.ID,
		Name:    testSong.Name,
		Artist:  testSong.Artist,
		Release: "16.07.2006",
		Text:    "Ooh baby, don't you know I suffer?\\nOoh baby, can\nyou hear me moan?\\nYou caught me under false pretenses\\nHow long\nbefore you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set\nmy soul alight",
		Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	var detailSong models.Song
	statusCode, detailSong, err = UpdateSong(t, testUpdateSong)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testUpdateSong, detailSong)

	// Проверяем получение информации по песне
	var songInfo models.SongInfoResponse
	statusCode, songInfo, err = GetSongByName(t, testSong)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testUpdateSong.Release, songInfo.Release)
	assert.Equal(t, testUpdateSong.Text, songInfo.Text)
	assert.Equal(t, testUpdateSong.Link, songInfo.Link)

	// Удаляем тестовые данные
	statusCode, err = DeleteSong(t, song)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}
