package dto

import "time"

type SongDbEntity struct {
	Song         string
	Song_Group   string
	Song_Text    string
	Link         string
	Release_Date time.Time
}

// POST
type SongPostRequest struct {
	Song  string `json:"song"`
	Group string `json:"group"`
}

type SongExtraInfoResponse struct {
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}

type SongCompleteInfo struct {
	Song        string
	Group       string
	Text        string
	Link        string
	ReleaseDate time.Time
}

type SongPostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PUT
type SongPutRequest struct {
	Song        string `json:"song"`
	Group       string `json:"group"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}

type SongPutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

//DELETE
type SongDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

//GET one

//GET multiple
type SongGetMultipleResponse struct {
}
