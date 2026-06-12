package database

import (
	"c2/server/models"
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type Db struct {
	DB *sql.DB
}

func NewDbConnection() *Db {

	db, err := sql.Open("sqlite", "bot.db")
	if err != nil {
		log.Println("Database connection error : ", err)
		return nil
	}
	return &Db{
		DB: db,
	}
}
func (db *Db) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS BotCreds(
						id TEXT PRIMARY KEY,
						registration_key TEXT NOT NULL
			  )`

	if _, err := db.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

func (db *Db) SaveToDB(models.BotCreds) {

}
