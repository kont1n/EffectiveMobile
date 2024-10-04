package tests

import (
	"EffectiveMobile/internal/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetSongList(t *testing.T) {
	t.Log("Get song list")
	var toDelete []models.SongResponse

	for i := 0; i < 10; i++ {
		t.Log("Add song")

		testSong := models.SongRequest{
			Name:   fmt.Sprintf("Song %v", i),
			Artist: "Muse",
		}

		statusCode, song, err := CreateSong(t, testSong)
		if err != nil {
			t.Log("Error creating song: ", err)
		}
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, statusCode)

		toDelete = append(toDelete, song)
	}

	statusCode, songsList, err := GetSongsList(t)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, len(toDelete), len(songsList.Songs))

	t.Log("songsList: ", songsList)

	// удаляем тестовые данные
	for _, song := range toDelete {
		statusCode, err = DeleteSong(t, song)
		if err != nil {
			t.Log("Error deleting song: ", err)
		}
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
	}

}
