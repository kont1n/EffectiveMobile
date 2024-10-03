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

func TestReadSong(t *testing.T) {
	t.Log("Read song")
	// создаем тестовые данные
	testSong := models.SongRequest{
		Name:   "Everybody's Fool",
		Artist: "Evanescence",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, song, err := CreateSong(t, testSong, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// проверяем что полученные данные совпадают с тестовыми
	var detailSong models.Song
	statusCode, detailSong, err = GetSongByID(t, song, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testSong.Name, detailSong.Name)
	assert.Equal(t, testSong.Artist, detailSong.Artist)

	// удаляем тестовые данные
	statusCode, err = DeleteSong(t, song, url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestReadSongWithEmptyID(t *testing.T) {
	t.Log("Read song with empty ID")

	getUrl := fmt.Sprintf("http://localhost:%s/api/song/%s", os.Getenv("WEBSERVER_PORT"), "")
	getResponse, err := http.Get(getUrl)
	defer getResponse.Body.Close()

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, getResponse.StatusCode)
}
