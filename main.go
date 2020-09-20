/*Package main is telegram bot @UC2_bot for https://t.me/UCChat
 */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//function main do main functions =)
func main() {
	//read bot token from env
	bottoken := os.Getenv("BOT_TOKEN")
	mysqlpass := os.Getenv("MYSQL_PASSWORD")
	if bottoken == "" {
		fmt.Fprintf(os.Stderr, "BOT TOKEN NOT FOUND!\n")
		os.Exit(1)
	}
	connectionstring := "root:" + mysqlpass + "@tcp(127.0.0.1:3306)/uc2botdatabase"
	db, err := sql.Open("mysql", connectionstring)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	//connect to telegram api
	bot, err := tgbotapi.NewBotAPI(bottoken)
	if err != nil {
		log.Panic(err)
	}
	t := time.Now()
	bot.Debug = true
	botname := bot.Self.UserName
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//read updates and do actions
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//welcome message for new users in group
		if update.Message.NewChatMembers != nil {
			for _, v := range *update.Message.NewChatMembers {
				chatname := update.Message.Chat.UserName
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+v.FirstName+"!\nДобро пожаловать в теплый ламповый чат: @"+chatname)
				bot.Send(msg)
			}
		}

		var re = regexp.MustCompile(`бекап|бэкап|рестор|backup|restore|ревакер|резервное|рекавер|восстанов`) //first encounter!
		now := time.Now()

		if re.MatchString(strings.ToLower(update.Message.Text)) && (now.After(t.Add(2 * time.Minute))) { //action about backup =)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Погугли =)")
			if update.Message.From.UserName == "angrypuffin" { //vovney exeption
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Погугли =)")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Только рековерить чангу надо по феншую, а не тупо ресторя всю вм. Вот тебе ссылочка от Вима на почитать https://www.veeam.com/wp-getting-best-availability-microsoft-exchange-veeam.html")
			}
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			t = time.Now()
		} else if update.Message.Text == "/help" || update.Message.Text == "/help@"+botname { // help action
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello world, Viva la @BrainFair!\n You can make me better: https://github.com/brainfair/uc2bot")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else if update.Message.Text == "/dbtest" { // dbtest action
			var answer string
			err := db.QueryRow("SELECT answer FROM QNA where question='q1'").Scan(&answer)
			// if there is an error inserting, handle it
			if err != nil {
				panic(err.Error())
			}
			// be careful deferring Queries if you are using transactions
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Наверно ответ: "+answer)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else if update.Message.Text == "/report" || update.Message.Text == "/report@"+botname {
			admins, err := bot.GetChatAdministrators(update.Message.Chat.ChatConfig()) //get chat admins
			if err != nil {
				panic(err.Error())
			}
			message := "Alarm summoning:"
			for _, adminuser := range admins {
				message += " @" + adminuser.User.UserName
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
	defer db.Close()
}
