package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"EffectiveMobile/internal/models"
)

func CreateSong(t *testing.T, testSong models.SongRequest, baseUrl string) (int, models.SongResponse, error) {
	t.Log("Calling the API to create a song")

	result := models.SongResponse{}
	body, err := json.Marshal(testSong)
	if err != nil {
		t.Errorf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf := bytes.NewBuffer(body)

	response, err := http.Post(baseUrl, "application/json", buf)
	if err != nil {
		t.Errorf("Error creating song: %v", err)
		return -1, result, err
	}
	defer response.Body.Close()

	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
		return -1, result, err
	}

	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	return response.StatusCode, result, err
}

func DeleteSong(t *testing.T, song models.SongResponse, baseUrl string) (int, error) {
	t.Log("Calling the API to delete a song")

	client := &http.Client{}
	deleteUrl := baseUrl + song.ID
	deleteRequest, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		t.Errorf("Error creating delete request: %v", err)
		return -1, err
	}

	deleteResponse, err := client.Do(deleteRequest)
	if err != nil {
		t.Errorf("Error deleting song: %v", err)
		return -1, err
	}
	defer deleteResponse.Body.Close()

	return deleteResponse.StatusCode, err
}

func GetSongByID(t *testing.T, song models.SongResponse, baseUrl string) (int, models.Song, error) {
	t.Log("Calling the API to get a song by ID")

	result := models.Song{}
	getUrl := baseUrl + song.ID
	response, err := http.Get(getUrl)
	if err != nil {
		t.Errorf("Error getting song: %v", err)
		return -1, result, err
	}
	defer response.Body.Close()

	var buf *bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
		return -1, result, err
	}

	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	return response.StatusCode, result, err
}

func UpdateSong(t *testing.T, testSong models.Song, baseUrl string) (statusCode int, result models.Song, err error) {
	t.Log("Calling the API to update a song")

	var body []byte
	var buf *bytes.Buffer
	var updateRequest *http.Request
	var updateResponse *http.Response

	body, err = json.Marshal(testSong)
	if err != nil {
		t.Errorf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf = bytes.NewBuffer(body)

	updateRequest, err = http.NewRequest(http.MethodPut, baseUrl, buf)
	if err != nil {
		t.Errorf("Error creating update request: %v", err)
		return -1, result, err
	}

	client := &http.Client{}
	updateResponse, err = client.Do(updateRequest)
	if err != nil {
		t.Errorf("Error updating song: %v", err)
		return -1, result, err
	}
	defer updateResponse.Body.Close()

	return updateResponse.StatusCode, result, err
}

func GetSongByName(t *testing.T, song models.SongRequest, baseUrl string) (statusCode int, result models.SongInfoResponse, err error) {
	t.Log("Calling the API to get a song by Name")

	var body []byte
	var buf *bytes.Buffer
	var getRequest *http.Request
	var getResponse *http.Response

	body, err = json.Marshal(song)
	if err != nil {
		t.Errorf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf = bytes.NewBuffer(body)

	getRequest, err = http.NewRequest(http.MethodGet, baseUrl, buf)
	if err != nil {
		t.Errorf("Error creating get request: %v", err)
		return -1, result, err
	}

	client := &http.Client{}
	getResponse, err = client.Do(getRequest)
	if err != nil {
		t.Errorf("Error getting song: %v", err)
		return -1, result, err
	}
	defer getResponse.Body.Close()

	return getResponse.StatusCode, result, err
}

func GetCouplet(t *testing.T, song models.SongResponse, coupletId string, baseUrl string) (statusCode int, result models.SongVerseResponse, err error) {
	t.Log("Calling the API to get a song couplet")

	getUrl := baseUrl + song.ID + "/verses/" + "/?coupletId=" + coupletId
	response, err := http.Get(getUrl)
	if err != nil {
		t.Errorf("Error getting song: %v", err)
		return -1, result, err
	}
	defer response.Body.Close()

	var buf *bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
		return -1, result, err
	}

	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	return response.StatusCode, result, err
}
