package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlexEr256/musicService/dto"
	"github.com/AlexEr256/musicService/repositories"
	"github.com/AlexEr256/musicService/utils"
	"github.com/gofiber/fiber/v2"
	"gitlab.com/metakeule/fmtdate"
	"go.uber.org/zap"
)

type ISongHandler interface {
	AddSong(c *fiber.Ctx) error
	UpdateSong(c *fiber.Ctx) error
	DeleteSong(c *fiber.Ctx) error
	GetSong(c *fiber.Ctx) error
	GetSongs(c *fiber.Ctx) error
}

type SongHandler struct {
	SongRepository repositories.ISongRepository
	Logger         *zap.Logger
}

func NewSongHandler(repository repositories.ISongRepository, logger *zap.Logger) ISongHandler {
	return &SongHandler{SongRepository: repository, Logger: logger}
}

func (h SongHandler) AddSong(c *fiber.Ctx) error {
	request := &dto.SongPostRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to parse request body - %s", err.Error()),
		})
	}
	h.Logger.Info("Add song to database",
		zap.String("Song name", request.Song),
		zap.String("Group", request.Group))

	if request.Song == "" || request.Group == "" {
		h.Logger.Error("Empty fields in request body", zap.String("Song name", request.Song), zap.String("Group", request.Group))
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
			Success: false,
			Message: "some fields of request body are empty",
		})
	}

	extraInfo, err := utils.SearchHook(request.Group, request.Song)

	if err != nil {
		h.Logger.Error("Failed to use api hook")
		return c.Status(fiber.StatusInternalServerError).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get info about song - %s", err.Error()),
		})
	}

	date, err := fmtdate.Parse("DD.MM.YYYY", extraInfo.ReleaseDate)
	if err != nil {
		h.Logger.Error("Failed parse date", zap.String("Date", extraInfo.ReleaseDate))
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to extract release date of the song - %s", err.Error()),
		})
	}

	fullSongInfo := &dto.SongCompleteInfo{
		Song:        request.Song,
		Group:       request.Group,
		Text:        extraInfo.Text,
		Link:        extraInfo.Link,
		ReleaseDate: date,
	}

	err = h.SongRepository.Add(fullSongInfo)

	if err != nil {
		h.Logger.Error("Failed to insert entity in database")
		return c.Status(fiber.StatusInternalServerError).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to save song entity in db - %s", err.Error()),
		})
	}
	h.Logger.Info("Insert entity in database")
	return c.Status(fiber.StatusCreated).JSON(&dto.SongPostResponse{
		Success: true,
	})
}
func (h SongHandler) UpdateSong(c *fiber.Ctx) error {
	song := c.Params("song")

	songInfo, err := h.SongRepository.GetOne(song)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get info about song - %s", err.Error()),
		})
	}

	request := &dto.SongPutRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to parse request body - %s", err.Error()),
		})
	}

	h.Logger.Info("Update song request body",
		zap.String("Song", request.Song),
		zap.String("Group", request.Group),
		zap.String("Text", request.Text),
		zap.String("Link", request.Link),
		zap.String("Release Date", request.ReleaseDate),
	)

	h.Logger.Info("Existed song info",
		zap.String("Song", songInfo.Song),
		zap.String("Group", songInfo.Song_Group),
		zap.String("Text", songInfo.Song_Text),
		zap.String("Link", songInfo.Link),
	)

	if request.Group != "" {
		songInfo.Song_Group = request.Group
	}
	if request.Link != "" {
		songInfo.Link = request.Link
	}
	if request.ReleaseDate != "" {
		date, err := fmtdate.Parse("DD.MM.YYYY", request.ReleaseDate)
		if err != nil {
			h.Logger.Error("Failed to parse release date", zap.String("Release Date", request.ReleaseDate))
			return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
				Success: false,
				Message: fmt.Sprintf("failed to extract release date of the song - %s", err.Error()),
			})
		}
		songInfo.Release_Date = date
	}
	if request.Text != "" {
		songInfo.Song_Text = request.Text
	}

	err = h.SongRepository.Update(songInfo)
	if err != nil {
		h.Logger.Error("Failed to update song in database")
		return c.Status(fiber.StatusInternalServerError).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update song info - %s", err.Error()),
		})
	}

	h.Logger.Info("Song was updated successfully")
	return c.Status(fiber.StatusOK).JSON(&dto.SongPostResponse{
		Success: true,
	})
}
func (h SongHandler) DeleteSong(c *fiber.Ctx) error {
	song := c.Params("song")

	h.Logger.Info("Delete song",
		zap.String("Song", song),
	)

	err := h.SongRepository.Delete(song)

	if err != nil {
		h.Logger.Error("Failed to delete song from database")
		return c.Status(fiber.StatusNotFound).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete song - %s", err.Error()),
		})
	}

	h.Logger.Info("Song was deleted from database")
	return c.Status(fiber.StatusOK).JSON(&dto.SongPostResponse{
		Success: true,
	})
}
func (h SongHandler) GetSong(c *fiber.Ctx) error {
	song := c.Params("song")
	page := c.Query("page", "0")

	pageValue, err := strconv.Atoi(page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert page query parameter to number - %s", err.Error()),
		})
	}
	if pageValue < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: "page query parameter should be positive number",
		})
	}

	songInfo, err := h.SongRepository.GetOne(song)

	h.Logger.Info("Existed song info",
		zap.Int("Page", pageValue),
		zap.String("Song", songInfo.Song),
		zap.String("Group", songInfo.Song_Group),
		zap.String("Text", songInfo.Song_Text),
		zap.String("Link", songInfo.Link),
	)

	if err != nil {
		h.Logger.Error("Failed to get song from database")
		return c.Status(fiber.StatusNotFound).JSON(&dto.SongGetResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get info about song - %s", err.Error()),
		})
	}

	songCouplets := strings.Split(songInfo.Song_Text, "\n")
	verses := make([]*dto.Verses, 0)

	for index, couplet := range songCouplets {
		if page != "0" && pageValue != index+1 {
			continue
		}
		verses = append(verses, &dto.Verses{Page: index + 1, Verse: couplet})
	}

	response := &dto.SongGetOneResponse{
		Song:        songInfo.Song,
		Group:       songInfo.Song_Group,
		Link:        songInfo.Link,
		Verses:      verses,
		ReleaseDate: songInfo.Release_Date,
	}

	return c.Status(fiber.StatusOK).JSON(&dto.SongGetResponse{
		SongInfo: response,
		Success:  true,
	})

}
func (h SongHandler) GetSongs(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageValue, err := strconv.Atoi(page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert page query parameter to number - %s", err.Error()),
		})
	}
	if pageValue < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: "page query parameter should be positive number",
		})
	}

	perPage := c.Query("perPage", "3")
	perPageValue, err := strconv.Atoi(perPage)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: fmt.Sprintf("failed to convert perPage query parameter to number - %s", err.Error()),
		})
	}
	if perPageValue < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetResponse{
			Success: false,
			Message: "perPage query parameter should be positive number",
		})
	}

	request := &dto.SongSearchRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to parse request body - %s", err.Error()),
		})
	}

	h.Logger.Info("Request info",
		zap.Int("Page", pageValue),
		zap.Int("perPage", perPageValue),
		zap.String("Song", request.Song),
		zap.String("Group", request.Group),
		zap.String("Text", request.Text),
		zap.String("Link", request.LinkFilter),
	)

	startDate := ""
	endDate := ""

	if request.ReleaseDateFilter != nil {
		if request.ReleaseDateFilter.End != "" {
			end, err := fmtdate.Parse("DD.MM.YYYY", request.ReleaseDateFilter.End)
			if err != nil {
				h.Logger.Error("Failed to convert end date to string", zap.String("End", request.ReleaseDateFilter.End))
				return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
					Success: false,
					Message: fmt.Sprintf("failed to extract end date filter of the song - %s", err.Error()),
				})
			}
			h.Logger.Info("End date",
				zap.String("End", request.ReleaseDateFilter.End),
			)
			endDateFormatted := fmtdate.FormatDate(end)
			endDate = endDateFormatted
		}

		if request.ReleaseDateFilter.Start != "" {
			start, err := fmtdate.Parse("DD.MM.YYYY", request.ReleaseDateFilter.Start)

			if err != nil {
				h.Logger.Error("Failed to convert start date to string", zap.String("Start", request.ReleaseDateFilter.Start))
				return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
					Success: false,
					Message: fmt.Sprintf("failed to extract start date filter of the song - %s", err.Error()),
				})
			}
			h.Logger.Info("Start date",
				zap.String("Start", request.ReleaseDateFilter.Start),
			)
			startDateFormatted := fmtdate.FormatDate(start)
			startDate = startDateFormatted
		}
	}

	songsInfo, err := h.SongRepository.GetMany(startDate, endDate, request.LinkFilter, request.Song, request.Group, request.Text)
	if err != nil {
		h.Logger.Error("Failed to extract info about songs")
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongGetManyResponse{
			Success: false,
			Message: fmt.Sprintf("failed to extract info about songs - %s", err.Error()),
		})
	}

	songs := make([]*dto.SongGetOneResponse, 0)
	for _, song := range songsInfo {
		songCouplets := strings.Split(song.Song_Text, "\n")
		verses := make([]*dto.Verses, 0)

		for index, couplet := range songCouplets {
			verses = append(verses, &dto.Verses{Page: index + 1, Verse: couplet})
		}

		songs = append(songs, &dto.SongGetOneResponse{
			Song:        song.Song,
			Group:       song.Song_Group,
			Link:        song.Link,
			Verses:      verses,
			ReleaseDate: song.Release_Date,
		})
	}

	total := len(songs) / perPageValue
	remains := len(songs) % perPageValue
	if remains != 0 {
		total = total + 1
	}
	offset := perPageValue * (pageValue - 1)

	if offset > len(songs) {
		offset = len(songs)
	}

	end := offset + perPageValue
	if end > len(songs) {
		end = len(songs)
	}

	h.Logger.Info("Extra info about request",
		zap.Int("Total", total),
		zap.Int("Offset", offset),
		zap.Int("Offset", perPageValue),
	)

	return c.Status(fiber.StatusOK).JSON(&dto.SongGetManyResponse{
		SongsInfo: songs[offset:end],
		TotalPage: total,
		Success:   true,
	})
}
