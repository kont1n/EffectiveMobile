package models

import "database/sql"

type Song struct {
	ID      string `json:"id"`
	Name    string `json:"song"`
	Artist  string `json:"group"`
	Release string `json:"releaseDate,omitempty"`
	Text    string `json:"text,omitempty"`
	Link    string `json:"link,omitempty"`
}

type SongDB struct {
	ID      string         `db:"id"`
	Name    string         `db:"song_name"`
	Artist  string         `db:"artist_name"`
	Release sql.NullTime   `db:"release_date"`
	Text    sql.NullString `db:"song_text"`
	Link    sql.NullString `db:"link"`
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

type SongInfoResponseDB struct {
	Release sql.NullTime
	Text    sql.NullString
	Link    sql.NullString
}

type SongsListResponse struct {
	Songs []Song `json:"songs"`
}

type SongsListResponseDB struct {
	Songs []SongDB `json:"songs"`
}

type SongVerseResponse struct {
	ID        string `json:"id"`
	CoupletId int    `json:"coupletId"`
	Couplet   string `json:"couplet"`
}

type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Error     string `json:"error"`
}

type SortOptions struct {
	Field string
	Order string
}

type PaginationOptions struct {
	Limit     string
	Offset    string
	PageToken string
}
