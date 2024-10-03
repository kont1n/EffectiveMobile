package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"EffectiveMobile/internal/models"
)

func TestCreateSong(t *testing.T) {
	t.Log("Create song")
	// Проверяем создание тестовых данных
	testSong := models.SongRequest{
		Name:   "Supermassive Black Hole",
		Artist: "Muse",
	}

	statusCode, song, err := CreateSong(t, testSong)
	if err != nil {
		t.Log("Error creating song: ", err)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// Удаляем тестовые данные
	statusCode, err = DeleteSong(t, song)
	if err != nil {
		t.Log("Error deleting song: ", err)
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

}

func TestCreateSongWithEmptyName(t *testing.T) {
	t.Log("Create song with empty name")

	testSong := models.SongRequest{
		Name:   "",
		Artist: "Muse",
	}

	statusCode, _, err := CreateSong(t, testSong)
	if err != nil {
		t.Log("Error creating song: ", err)
	}

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
