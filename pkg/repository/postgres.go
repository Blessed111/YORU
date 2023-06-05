package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

const teamtable = "Teams"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Sslmode  string
}

// func NewPostgresDb(cfg Config) (*sqlx.DB, error) {
// 	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

func NewPostgresDb1() (*sqlx.DB, error) {
	username := "postgres"
	password := "55313104"
	host := "localhost"
	port := "5432"
	dbname := "YORU"

	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, dbname))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	err = createTables(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sqlx.DB) error {
	path := "../schema/"
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Printf("error when read dir tables: %v", err)
		return err
	}
	for _, v := range dir {
		query, err := os.ReadFile(path + v.Name())
		if err != nil {
			log.Printf("error when read file tables: %v", err)
			return err
		}

		_, err = db.Exec(string(query))
		if err != nil {
			log.Printf("error when create table: %v", err)
			return err
		}
	}

	return nil
}
