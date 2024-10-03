package tests

import (
	"net/http"
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

	statusCode, song, err := CreateSong(t, testSong)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	// проверяем удаление
	statusCode, err = DeleteSong(t, song)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}
