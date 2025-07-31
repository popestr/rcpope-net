package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	secretName = "rcpopenet_verceldb"
)

var (
	connection PostgresConnection
)

type PostgresConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
}

func init() {
	secret, err := GetSecret(secretName)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(secret), &connection)
	if err != nil {
		log.Fatal(err)
	}
}

func MustGet(secretName string) string {
	val, ok := os.LookupEnv(secretName)
	if !ok {
		panic("missing required environment variable: " + secretName)
	}
	return val
}

func (c PostgresConnection) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", c.Host, c.Port, c.User, c.Password, c.DbName)
}

func DB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connection.GetConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
