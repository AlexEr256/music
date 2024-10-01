package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type SqlConnection struct {
	Db *sqlx.DB
}

func NewConnection(connection string) (SqlConnection, error) {
	db, err := sqlx.Connect("pgx", connection)

	if err != nil {
		return SqlConnection{}, err
	}

	if err = db.Ping(); err != nil {
		return SqlConnection{}, err
	}

	return SqlConnection{Db: db}, nil
}
