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

func TestCreateSong(t *testing.T) {
	t.Log("Create song")

	testSong := models.SongRequest{
		Name:   "Supermassive Black Hole",
		Artist: "Muse",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, _, err := CreateSong(t, testSong, url)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)
}

func TestCreateSongWithEmptyName(t *testing.T) {
	t.Log("Create song with empty name")

	testSong := models.SongRequest{
		Name:   "",
		Artist: "Muse",
	}

	url := fmt.Sprintf("http://localhost:%s/api/song", os.Getenv("WEBSERVER_PORT"))
	statusCode, _, err := CreateSong(t, testSong, url)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
