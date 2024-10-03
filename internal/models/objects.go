package models

type Song struct {
	ID      string `json:"id"`
	Name    string `json:"song"`
	Artist  string `json:"group"`
	Release string `json:"release_date,omitempty"`
	Text    string `json:"text,omitempty"`
	Link    string `json:"link,omitempty"`
}

type SongRequest struct {
	Name   string `json:"song"`
	Artist string `json:"group"`
}

type SongResponse struct {
	ID string `json:"id"`
}

type SongInfoResponse struct {
	Release string `json:"releaseDate"`
	Text    string `json:"text"`
	Link    string `json:"link"`
}

type SongsListResponse struct {
	Songs []Song `json:"songs"`
}

type SongVerseResponse struct {
	ID        string `json:"id"`
	Name      string `json:"song"`
	Artist    string `json:"group"`
	CoupletId int    `json:"coupletId"`
	Couplet   string `json:"couplet"`
}

type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Error     string `json:"error"`
}
