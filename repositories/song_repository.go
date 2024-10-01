package repositories

import (
	"fmt"

	"github.com/AlexEr256/musicService/dto"
	"github.com/jmoiron/sqlx"
)

type ISongRepository interface {
	Add(c *dto.SongCompleteInfo) error
	Update(c *dto.SongDbEntity) error
	Delete() (*dto.SongDeleteResponse, error)
	GetOne(song string) (*dto.SongDbEntity, error)
	GetMany() (*dto.SongGetMultipleResponse, error)
}

type SongRepository struct {
	Db *sqlx.DB
}

func NewSongRepository(db *sqlx.DB) ISongRepository {
	return &SongRepository{Db: db}
}

func (r SongRepository) Add(songRequest *dto.SongCompleteInfo) error {
	query := `INSERT INTO
				songs(song_group, song, song_text, link, release_date)
			VALUES
				(:group, :song, :text, :link, :releasedate);`

	_, err := r.Db.NamedExec(query, songRequest)
	if err != nil {
		return err
	}

	return nil
}

func (r SongRepository) Update(songRequest *dto.SongDbEntity) error {
	query := `UPDATE songs
				SET song_group=$1, song_text=$2, link=$3, release_date=$4
			WHERE
				song=$5`

	_, err := r.Db.Query(query, songRequest.Song_Group, songRequest.Song_Text, songRequest.Link, songRequest.Release_Date, songRequest.Song)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (r SongRepository) Delete() (*dto.SongDeleteResponse, error) {
	return nil, nil
}

func (r SongRepository) GetOne(song string) (*dto.SongDbEntity, error) {
	resp := dto.SongDbEntity{}
	err := r.Db.Get(&resp, "SELECT * FROM songs WHERE song=$1", song)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r SongRepository) GetMany() (*dto.SongGetMultipleResponse, error) {
	return nil, nil
}
