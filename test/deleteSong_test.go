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

func TestDeleteSong(t *testing.T) {
	t.Log("Delete song")
	// создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Dream On",
		Artist: "Aerosmith",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, song, err := CreateSong(t, testSong, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// проверяем удаление
	statusCode, err = DeleteSong(t, song, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestDeleteSongWithEmptyID(t *testing.T) {
	t.Log("Delete song with empty ID")

	song := models.SongResponse{
		ID: "",
	}
	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, err := DeleteSong(t, song, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
