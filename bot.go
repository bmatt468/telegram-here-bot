package herebot

import (
    "database/sql"
	"log"
    "github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
    "strings"
)

type HereBot struct {
    API *tgbotapi.BotAPI
    dbPass string
}

func Connect(tkn string, debug bool, dbPass string) (*SmawkBot, error) {
    // Call the Telegram API wrapper and authenticate our Bot
    bot, err := tgbotapi.NewBotAPI(tkn)

    // Check to see if there were any errors with our bot and fail
    // if there were
    if err != nil {
        log.Fatal(err)
    }

    if (debug) {
        // Print confirmation
        log.Printf("Authorized on account %s", bot.Self.UserName)
    }

    // Set our bot to either be in debug mode (everything gets put out to the console)
    // or non debug mode (everything is silent)
    bot.Debug = debug

    // Create the SmawkBot
    hbot := &HereBot {
        API: bot,
        dbPass: dbPass,
    }

    // Return our bot back to the caller
    return hbot, err
}

// OpenWebhook opens up a webhook without attaching a self signed certificate
func (bot *HereBot) OpenWebhook(url string) {
    _, err := bot.API.SetWebhook(tgbotapi.NewWebhook(url))
    if err != nil {
        log.Fatal(err)
    }
}

func (bot *HereBot) Listen(token string) <-chan tgbotapi.Update {
    updates := bot.API.ListenForWebhook(token)
    return updates
}

func (bot *HereBot) ParseAndExecuteUpdate(update tgbotapi.Update) {
    if update.Message != nil {
        switch cmd := strings.Split(update.Message.Text, " "); strings.Replace(cmd[0],"@everyone_here_bot","",-1) {
        case "/start":
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "HereBot is now active.")
            bot.API.Send(msg)

        case "/register":
            db, err := ConnectDB(bot.dbPass)
            if err != nil {
                log.Fatal(err)
            }
            defer db.Close()

            if len(cmd) > 1 {
                msg_string := "Correct Usage: /register"
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                bot.API.Send(msg)
            } else {
                result, err := db.Query("INSERT INTO users(username,chat_id,flag_active) VALUES(?,?,1)","@"+update.Message.From.UserName,update.Message.Chat.ID)
                if err != nil {
                    log.Fatal(err)
                }
                defer result.Close()

                msg_string := "@"+update.Message.From.UserName+" has been registered."
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                bot.API.Send(msg)
            }

        case "/deregister":
            db, err := ConnectDB(bot.dbPass)
            if err != nil {
                log.Fatal(err)
            }
            defer db.Close()

            if len(cmd) > 1 {
                msg_string := "Correct Usage: /deregister"
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                bot.Send(msg)
            } else {
                result, err := db.Query("UPDATE users SET flag_active=0 WHERE chat_id=? AND username=? ",update.Message.Chat.ID,"@"+update.Message.From.UserName)
                if err != nil {
                    log.Fatal(err)
                }
                defer result.Close()

                msg_string := "@"+update.Message.From.UserName+" has been deregistered."
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                bot.API.Send(msg)
            }

        case "/all", "/here":
            db, err := ConnectDB(bot.dbPass)
            if err != nil {
                log.Fatal(err)
            }
            defer db.Close()

            users, err := db.Query("SELECT username FROM users u WHERE chat_id=? AND flag_active=1 GROUP BY u.username, u.chat_id",update.Message.Chat.ID)
            if err != nil {
                log.Fatal(err)
            }
            defer users.Close()

            msg_string := ""

            for users.Next() {
                var username string
                if err := users.Scan(&username); err != nil {
                    log.Fatal(err)
                }
                msg_string += " " + username
            }
            if err := users.Err(); err != nil {
                log.Fatal(err)
            }

            if msg_string == "" {
                msg_string = "No users registered."
            }

            msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
            bot.API.Send(msg)
        }
    }
}

func ConnectDB(password string) (*sql.DB, error) {
    cfg := &mysql.Config {
        User: "hello-bot",
        Passwd: password,
        Net: "tcp",
        Addr: "107.170.45.12:3306",
        DBName: "hello-bot",
    }
    return sql.Open("mysql", cfg.FormatDSN())
}

