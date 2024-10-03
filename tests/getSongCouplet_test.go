package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"EffectiveMobile/internal/models"
)

func TestGetSongCouplet(t *testing.T) {
	t.Log("Get song couplet")
	// Создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Everybody's Fool",
		Artist: "Evanescence",
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

	// Проверяем получение куплета
	couplets := strings.Split(testUpdateSong.Text, "\\n\\n")
	response := models.SongVerseResponse{}
	statusCode, response, err = GetCouplet(t, song, "0")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, couplets[0], response.Couplet)

	statusCode, response, err = GetCouplet(t, song, "1")
	require.NoError(t, err)
	assert.Equal(t, couplets[1], response.Couplet)

	// Удаляем тестовые данные
	statusCode, err = DeleteSong(t, song)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}
