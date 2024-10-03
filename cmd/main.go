package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlexEr256/musicService/database"
	_ "github.com/AlexEr256/musicService/database"
	"github.com/AlexEr256/musicService/handlers"
	"github.com/AlexEr256/musicService/repositories"
	"github.com/AlexEr256/musicService/utils"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading dotenv file")
	}

	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	db := os.Getenv("PG_DB")

	validationErrors := utils.ValidateEnvParams(user, pass, host, port, db)
	if len(validationErrors) != 0 {
		log.Fatal("Check input parameters ", validationErrors)
	}

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)

	pg, err := database.NewConnection(connection)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
		return
	}

	m, err := migrate.New(
		"file://migrations",
		connection,
	)

	if err != nil {
		log.Fatal("Failed to prepare migrations ", err)
	}

	if err = m.Up(); err != migrate.ErrNoChange && err != nil {
		log.Fatal("Failed to execute migrations ", err)
	}

	r := repositories.NewSongRepository(pg.Db)
	h := handlers.NewSongHandler(r)

	app := fiber.New()

	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger_shop.json",
		Path:     "docs",
		Title:    "Swagger API Docs",
	}
	app.Use(swagger.New(cfg))

	app.Post("/songs/search", h.GetSongs)

	app.Post("/songs", h.AddSong)
	app.Put("/songs/:song", h.UpdateSong)
	app.Delete("/songs/:song", h.DeleteSong)
	app.Get("/songs/:song", h.GetSong)

	app.Listen(":3000")
}
