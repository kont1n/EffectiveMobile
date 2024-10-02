package models

type SongDetail struct {
	ID      string `json:"id"`
	Name    string `json:"song"`
	Artist  string `json:"group"`
	Release string `json:"release_date,omitempty"`
	Text    string `json:"text,omitempty"`
	Link    string `json:"link,omitempty"`
}

type SongPostRequest struct {
	Name   string `json:"song"`
	Artist string `json:"group"`
}

type SongPostResponse struct {
	ID string `json:"id"`
}

type SongsGetResponse struct {
	Songs []SongDetail `json:"songs"`
}

type SongVerseResponse struct {
	ID        string `json:"id"`
	Name      string `json:"song"`
	Artist    string `json:"group"`
	CoupletId int    `json:"couplet_id"`
	Couplet   string `json:"couplet"`
}

type ErrorResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
}
