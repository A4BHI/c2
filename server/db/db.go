package database

import (
	"c2/server/models"
	"database/sql"
	"log"
)

func StartDB() {
	db, err := sql.Open("sqlite", "bot.db")
	if err != nil {
		log.Println("Database connection error : ", err)
		return
	}

}

func SaveToDB(models.BotCreds) {

}
