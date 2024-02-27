package postgres

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() (*sqlx.DB, error) {
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	dataSourceName := "user=" + pgUser + " password=" + pgPassword + " dbname=" + pgDB + " host=" + pgHost + " port=" + pgPort + " sslmode=disable"

	postgres, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return postgres, nil
}
