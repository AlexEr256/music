package handlers

import (
	"fmt"

	"github.com/AlexEr256/musicService/dto"
	"github.com/AlexEr256/musicService/repositories"
	"github.com/AlexEr256/musicService/utils"
	"github.com/gofiber/fiber/v2"
	"gitlab.com/metakeule/fmtdate"
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
}

func NewSongHandler(repository repositories.ISongRepository) ISongHandler {
	return &SongHandler{SongRepository: repository}
}

func (h SongHandler) AddSong(c *fiber.Ctx) error {
	request := &dto.SongPostRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to parse request body - %s", err.Error()),
		})
	}

	if request.Song == "" || request.Group == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
			Success: false,
			Message: "some fields of request body are empty",
		})
	}

	extraInfo, err := utils.DumbSearchHook(request.Group, request.Song)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get info about song - %s", err.Error()),
		})
	}

	date, err := fmtdate.Parse("DD.MM.YYYY", extraInfo.ReleaseDate)
	if err != nil {
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
		return c.Status(fiber.StatusInternalServerError).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to save song entity in db - %s", err.Error()),
		})
	}

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

	if request.Group != "" {
		songInfo.Song_Group = request.Group
	}
	if request.Link != "" {
		songInfo.Link = request.Link
	}
	if request.ReleaseDate != "" {
		date, err := fmtdate.Parse("DD.MM.YYYY", request.ReleaseDate)
		if err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(&dto.SongPostResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update song info - %s", err.Error()),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&dto.SongPostResponse{
		Success: true,
	})
}
func (h SongHandler) DeleteSong(c *fiber.Ctx) error { return nil }
func (h SongHandler) GetSong(c *fiber.Ctx) error    { return nil }
func (h SongHandler) GetSongs(c *fiber.Ctx) error   { return nil }
