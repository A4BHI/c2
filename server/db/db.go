package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type Db struct {
	Conn *sql.DB
}

func NewDbConnection() *Db {

	db, err := sql.Open("sqlite", "bot.db")
	if err != nil {
		log.Println("Database connection error : ", err)
		return nil
	}
	return &Db{
		Conn: db,
	}
}
func (db *Db) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS BotCreds(
						id TEXT PRIMARY KEY,
						registration_key TEXT NOT NULL
			  )`

	if _, err := db.Conn.Exec(query); err != nil {
		return err
	}

	return nil
}

func (db *Db) SavetoDB(query string, args ...any) error {
	if _, err := db.Conn.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (db *Db) SearchAgentID(id string) bool {
	query := `SELECT CASE
    WHEN EXISTS (
        SELECT 1
        FROM BotCreds
        WHERE id = ?
    )
    THEN 1
    ELSE 0
END AS data_exists;`
	var exists bool
	err := db.Conn.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false
	}

	return true
}

func (db *Db) GetFromDB(query string, args ...any) any {

}
