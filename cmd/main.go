package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
	if user == "" {
		log.Fatal("PG_USER value is empty. Check input parameters")
	}

	pass := os.Getenv("PG_PASSWORD")
	if pass == "" {
		log.Fatal("PG_PASSWORD value is empty. Check input parameters")
	}

	host := os.Getenv("PG_HOST")
	if host == "" {
		log.Fatal("PG_HOST value is empty. Check input parameters")
	}

	port := os.Getenv("PG_PORT")
	if port == "" {
		log.Fatal("PG_PORT value is empty. Check input parameters")
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal("PG_PORT must be valid number. Get ", port)
	}

	database := os.Getenv("PG_DB")
	if database == "" {
		log.Fatal("PG_DB value is empty. Check input parameters")
	}

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("Failed to establish db connection ", err)
	}

	m, err := migrate.New(
		"file://migrations",
		conn,
	)

	if err != nil {
		log.Fatal("Failed to prepare migrations ", err)
	}

	if err = m.Up(); err != migrate.ErrNoChange && err != nil {
		log.Fatal("Failed to execute migrations ", err)
	}

	fmt.Println(db)
}
