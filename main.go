package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bottoken := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(bottoken)
	if err != nil {
		log.Panic(err)
	}
	t := time.Now()
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.NewChatMembers != nil {
			for _, v := range *update.Message.NewChatMembers {
				chatname := update.Message.Chat.UserName
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+v.FirstName+"!\nДобро пожаловать в теплый ламповый чат: @"+chatname)
				bot.Send(msg)
			}
		}

		var re = regexp.MustCompile(`бекап|бэкап|рестор|backup|restore|ревакер|резервное|рекавер|восстанов`)
		now := time.Now()
		fmt.Println(now)
		if re.MatchString(strings.ToLower(update.Message.Text)) && (now.After(t.Add(2 * time.Minute))) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Погугли =)")
			if update.Message.From.UserName == "angrypuffin" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Погугли =)")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Только рековерить чангу надо по феншую, а не тупо ресторя всю вм. Вот тебе ссылочка от Вима на почитать https://www.veeam.com/wp-getting-best-availability-microsoft-exchange-veeam.html")
			}
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			t = time.Now()
		} else if update.Message.Text == "/help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello world, Viva la @BrainFair!\n You can make me better: https://github.com/brainfair/uc2bot")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
