package repositories

import (
	"fmt"
	"time"

	"github.com/AlexEr256/musicService/dto"
	"github.com/jmoiron/sqlx"
)

type ISongRepository interface {
	Add(c *dto.SongCompleteInfo) error
	Update(c *dto.SongDbEntity) error
	Delete(song string) error
	GetOne(song string) (*dto.SongDbEntity, error)
	GetMany(startDate, endDate, link, song, group, text string) ([]*dto.SongDbEntity, error)
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

func (r SongRepository) Delete(song string) error {
	query := `DELETE FROM songs
	WHERE song=$1`

	result, err := r.Db.Exec(query, song)
	if err != nil {
		return err
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if deletedRows == 0 {
		return fmt.Errorf("no entities were deleted from database")
	}

	return nil
}

func (r SongRepository) GetOne(song string) (*dto.SongDbEntity, error) {
	resp := dto.SongDbEntity{}
	err := r.Db.Get(&resp, "SELECT * FROM songs WHERE song=$1", song)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r SongRepository) GetMany(startDate, endDate, link, song, group, text string) ([]*dto.SongDbEntity, error) {
	songs := make([]*dto.SongDbEntity, 0)
	finalStatement := ""
	sqlStatement := `SELECT * FROM songs`
	query := constructQuery(startDate, endDate, link, song, group, text)
	if query != "" {
		finalStatement = sqlStatement + " WHERE " + query
	} else if query == "" {
		finalStatement = sqlStatement + ";"
	}

	fmt.Println(finalStatement)
	rows, err := r.Db.Query(finalStatement)
	if err != nil {
		return nil, fmt.Errorf("failed to select songs - %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		song, group, text, link := "", "", "", ""
		var release time.Time
		var err = rows.Scan(&song, &group, &text, &link, &release)
		if err != nil {
			return songs, fmt.Errorf("failed to handle query result - %w", err)
		}
		songs = append(songs, &dto.SongDbEntity{Song: song, Song_Group: group, Song_Text: text, Link: link, Release_Date: release})
	}

	err = rows.Err()
	if err != nil {
		return songs, fmt.Errorf("failed to handle query result - %w", err)
	}

	return songs, nil

}

func constructQuery(startDate, endDate, link, song, group, text string) string {
	query := ""

	if query != "" && startDate != "" && endDate != "" {
		query += " AND release_date between " + fmt.Sprintf("%s%s%s", "'", startDate, "'") + " AND " + fmt.Sprintf("%s%s%s", "'", endDate, "'")
	} else if startDate != "" && endDate != "" && query == "" {
		query += "( release_date between " + fmt.Sprintf("%s%s%s", "'", startDate, "'") + " AND " + fmt.Sprintf("%s%s%s", "'", endDate, "'")
	}

	if query != "" && link != "" {
		query += " AND link LIKE " + fmt.Sprintf("%s%s%s", "'%", link, "%'")
	} else if link != "" && query == "" {
		query += "( link LIKE " + fmt.Sprintf("%s%s%s", "'%", link, "%'")
	}

	if query != "" && song != "" {
		query += " AND song LIKE " + fmt.Sprintf("%s%s%s", "'%", song, "%'")
	} else if song != "" && query == "" {
		query += "( song LIKE " + fmt.Sprintf("%s%s%s", "'%", song, "%'")
	}

	if query != "" && group != "" {
		query += " AND song_group LIKE " + fmt.Sprintf("%s%s%s", "'%", group, "%'")
	} else if group != "" && query == "" {
		query += "( song_group LIKE " + fmt.Sprintf("%s%s%s", "'%", group, "%'")
	}

	if query != "" && text != "" {
		query += " AND song_text LIKE " + fmt.Sprintf("%s%s%s", "'%", text, "%'")
	} else if text != "" && query == "" {
		query += "( song_text LIKE " + fmt.Sprintf("%s%s%s", "'%", text, "%'")
	}

	if query != "" {
		query += " );"
	}

	return query
}
