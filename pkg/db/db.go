package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	port = "5432"
)

var (
	host     = MustGet("POSTGRES_HOST")
	user     = MustGet("POSTGRES_USER")
	password = MustGet("POSTGRES_PASSWORD")
	dbname   = MustGet("POSTGRES_DATABASE")
)

func MustGet(secretName string) string {
	val, ok := os.LookupEnv(secretName)
	if !ok {
		panic("missing required environment variable: " + secretName)
	}
	return val
}

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
}

func GetDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", GetConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
