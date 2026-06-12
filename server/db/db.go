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

	query := `CREATE TABLE IF NOT EXISTS BotCreds(
						id TEXT PRIMARY KEY,
						registration_key TEXT NOT NULL
			  )`

	if _, err = db.Exec(query); err != nil {
		log.Println("Failed to execute query : ", err)
		return nil
	}

	return &Db{
		DB: db,
	}
}

func SaveToDB(models.BotCreds) {

}
