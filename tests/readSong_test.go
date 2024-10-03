package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"EffectiveMobile/internal/models"
)

func TestReadSong(t *testing.T) {
	t.Log("Read song")
	// создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Конь",
		Artist: "ЛЮБЭ",
	}

	statusCode, song, err := CreateSong(t, testSong)
	if err != nil {
		t.Log("Error creating song: ", err)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// проверяем что полученные данные совпадают с тестовыми
	var detailSong models.Song
	statusCode, detailSong, err = GetSongByID(t, song)
	if err != nil {
		t.Log("Error getting song: ", err)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testSong.Name, detailSong.Name)
	assert.Equal(t, testSong.Artist, detailSong.Artist)

	// удаляем тестовые данные
	statusCode, err = DeleteSong(t, song)
	if err != nil {
		t.Log("Error deleting song: ", err)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}
