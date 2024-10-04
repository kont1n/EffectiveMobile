package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"EffectiveMobile/internal/models"
)

const (
	baseUrl = "http://localhost:8080/api/song"
)

func CreateSong(t *testing.T, testSong models.SongRequest) (int, models.SongResponse, error) {
	t.Log("Calling the API to create a song")

	result := models.SongResponse{}
	body, err := json.Marshal(testSong)
	if err != nil {
		t.Logf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf := bytes.NewBuffer(body)

	t.Log("Sending request to ", baseUrl)
	t.Log("Sending request body", string(body))
	response, err := http.Post(baseUrl, "application/json", buf)
	if err != nil {
		t.Logf("Error creating song: %v", err)
		return -1, result, err
	}
	defer response.Body.Close()

	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		t.Logf("Error reading response body: %v", err)
		return -1, result, err
	}

	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}
	t.Log("Created song: ", result.ID)
	t.Log("Status code: ", response.StatusCode)

	return response.StatusCode, result, err
}

func DeleteSong(t *testing.T, song models.SongResponse) (int, error) {
	t.Log("Calling the API to delete a song")

	client := &http.Client{}
	deleteUrl := baseUrl + "/" + song.ID
	t.Log("Sending request to ", deleteUrl)
	deleteRequest, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		t.Logf("Error creating delete request: %v", err)
		return -1, err
	}

	deleteResponse, err := client.Do(deleteRequest)
	if err != nil {
		t.Logf("Error deleting song: %v", err)
		return -1, err
	}
	defer deleteResponse.Body.Close()
	t.Log("Deleted song: ", song.ID)
	t.Log("Status code: ", deleteResponse.StatusCode)

	return deleteResponse.StatusCode, err
}

func GetSongByID(t *testing.T, song models.SongResponse) (statusCode int, songDetail models.Song, err error) {
	t.Log("Calling the API to get a song by ID")

	var response *http.Response

	getUrl := baseUrl + "/" + song.ID
	t.Log("Sending request to ", getUrl)

	response, err = http.Get(getUrl)
	if err != nil {
		t.Logf("Error getting song: %v", err)
		return -1, songDetail, err
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Logf("Error reading response body: %v", err)
		return -1, songDetail, err
	}

	t.Log("Response: ", string(responseData))

	if err = json.Unmarshal(responseData, &songDetail); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, songDetail, err
	}

	return response.StatusCode, songDetail, err
}

func UpdateSong(t *testing.T, testSong models.Song) (statusCode int, result models.Song, err error) {
	t.Log("Calling the API to update a song")

	var body []byte
	var buf *bytes.Buffer
	var updateRequest *http.Request
	var updateResponse *http.Response

	body, err = json.Marshal(testSong)
	if err != nil {
		t.Logf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf = bytes.NewBuffer(body)

	updateRequest, err = http.NewRequest(http.MethodPut, baseUrl, buf)
	if err != nil {
		t.Logf("Error creating update request: %v", err)
		return -1, result, err
	}

	client := &http.Client{}
	t.Log("Sending request to ", baseUrl)
	t.Log("Sending request body", string(body))
	updateResponse, err = client.Do(updateRequest)
	if err != nil {
		t.Logf("Error updating song: %v", err)
		return -1, result, err
	}
	defer updateResponse.Body.Close()

	_, err = buf.ReadFrom(updateResponse.Body)
	if err != nil {
		t.Logf("Error reading response body: %v", err)
		return -1, result, err
	}
	t.Log("Response body: ", string(buf.Bytes()))

	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	t.Log("Updated song: ", result.ID)

	return updateResponse.StatusCode, result, err
}

func GetSongByName(t *testing.T, song models.SongRequest) (statusCode int, result models.SongInfoResponse, err error) {
	t.Log("Calling the API to get a song by Name")

	var body []byte
	var buf *bytes.Buffer
	var getRequest *http.Request
	var getResponse *http.Response

	body, err = json.Marshal(song)
	if err != nil {
		t.Logf("Error marshalling song: %v", err)
		return -1, result, err
	}
	buf = bytes.NewBuffer(body)

	getUrl := baseUrl + "/info"
	t.Log("Sending request to ", getUrl)
	t.Log("Sending request body", string(body))

	getRequest, err = http.NewRequest(http.MethodGet, getUrl, buf)
	if err != nil {
		t.Logf("Error creating get request: %v", err)
		return -1, result, err
	}
	getResponse, err = http.DefaultClient.Do(getRequest)
	if err != nil {
		t.Logf("Error getting song: %v", err)
		return -1, result, err
	}
	defer getResponse.Body.Close()

	responseData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		t.Logf("Error reading response body: %v", err)
		return -1, result, err
	}

	t.Log("Response: ", string(responseData))

	if err = json.Unmarshal(responseData, &result); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	return getResponse.StatusCode, result, err
}

func GetCouplet(t *testing.T, song models.SongResponse, coupletId string) (statusCode int, result models.SongVerseResponse, err error) {
	t.Log("Calling the API to get a song couplet")

	var getRequest *http.Request
	var getResponse *http.Response

	getUrl := baseUrl + "/" + song.ID + "/couplet" + "?id=" + coupletId

	t.Log("Sending request to ", getUrl)
	client := &http.Client{}
	getRequest, err = http.NewRequest(http.MethodGet, getUrl, nil)
	if err != nil {
		t.Logf("Error creating get request: %v", err)
		return -1, result, err
	}
	getRequest.Close = true

	getRequest.Header.Set("Content-Type", "application/json")
	getResponse, err = client.Do(getRequest)
	if err != nil {
		t.Logf("Error getting song: %v", err)
		return -1, result, err
	}
	defer getResponse.Body.Close()

	responseData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	t.Log("Response: ", string(responseData))

	if err = json.Unmarshal(responseData, &result); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, result, err
	}

	return getResponse.StatusCode, result, err
}

func GetSongsList(t *testing.T) (statusCode int, songsList models.SongsListResponse, err error) {
	t.Log("Calling the API to get a songs list")

	var response *http.Response

	getUrl := baseUrl + "s"
	t.Log("Sending request to ", getUrl)

	response, err = http.Get(getUrl)
	if err != nil {
		t.Logf("Error getting song: %v", err)
		return -1, songsList, err
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Logf("Error reading response body: %v", err)
		return -1, songsList, err
	}

	t.Log("Response: ", string(responseData))

	if err = json.Unmarshal(responseData, &songsList); err != nil {
		t.Logf("Error unmarshalling response body: %v", err)
		return -1, songsList, err
	}

	return response.StatusCode, songsList, err
}
