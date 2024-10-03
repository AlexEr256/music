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

type SongGetResponse struct {
	SongInfo *SongGetOneResponse `json:"songInfo"`
	Success  bool                `json:"success"`
	Message  string              `json:"message"`
}

type Verses struct {
	Page  int    `json:"page"`
	Verse string `json:"verse"`
}

type SongGetOneResponse struct {
	Song        string    `json:"song"`
	Group       string    `json:"group"`
	Link        string    `json:"link"`
	Verses      []*Verses `json:"verses"`
	ReleaseDate time.Time `json:"releaseDate"`
}

//GET multiple

type SongGetManyResponse struct {
	SongsInfo []*SongGetOneResponse `json:"songsInfo"`
	TotalPage int                   `json:"totalPages"`
	Success   bool                  `json:"success"`
	Message   string                `json:"message"`
}

type SongSearchRequest struct {
	ReleaseDateFilter *ReleaseDateFilter `json:"release"`
	LinkFilter        string             `json:"link"`
	Song              string             `json:"song"`
	Group             string             `json:"group"`
	Text              string             `json:"text"`
}

type ReleaseDateFilter struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
