package main

import (
    "database/sql"
	"log"
    "github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
    "strings"
)

const (
    bot_token = "293545528:AAHM-jN6D4Y84sNMlvhYHsf6CrblpJ0-nAI"
    bot_url = "https://www.benjaminrmatthews.com:88"
    bot_url_token = "/Dk39s0dk3S5PO12"
    bot_cert = "/etc/letsencrypt/live/benjaminrmatthews.com/fullchain.pem"
    bot_key = "/etc/letsencrypt/live/benjaminrmatthews.com/privkey.pem"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(bot_token)
    if err != nil {
        log.Fatal(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    _, err = bot.SetWebhook(tgbotapi.NewWebhook(bot_url+bot_url_token))
    if err != nil {
        log.Fatal(err)
    }

    updates := bot.ListenForWebhook(bot_url_token)
    go http.ListenAndServeTLS("0.0.0.0:88",bot_cert,bot_key,nil)

    for update := range updates {
        if update.Message != nil {
            switch cmd := strings.Split(update.Message.Text, " "); strings.Replace(cmd[0],"@everyone_here_bot","",-1) {
            case "/start":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "HereBot is now active.")
                bot.Send(msg)

            case "/register":
                log.Printf("%+v\n", cmd)

            case "/deregister":
                // Connect to our database
                db, err := ConnectDB()
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()

                if len(cmd) == 1 || len(cmd) > 2 {
                    // Wrong Usage
                    msg_string := "Correct Usage: /deregister @username"
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                    bot.Send(msg)

                } else if cmd[1] != "@"+update.Message.From.UserName {
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You may not deregister another user.")
                    bot.Send(msg)
                } else if len(cmd) == 2 {
                    result, err := db.Query("UPDATE users SET flag_active=0 WHERE chat_id=? AND username=? ",update.Message.Chat.ID,cmd[1])
                    if err != nil {
                        log.Fatal(err)
                    }
                    defer result.Close()

                    msg_string := cmd[1]+" has been deregistered."
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                    bot.Send(msg)
                }

            case "/all", "/here":
                db, err := ConnectDB()
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()

                users, err := db.Query("SELECT username FROM users u WHERE chat_id=? AND flag_active=1",update.Message.Chat.ID)
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

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
                bot.Send(msg)
            }
        }
    }
}

func ConnectDB() (*sql.DB, error) {
    cfg := &mysql.Config {
        User: "hello-bot",
        Passwd: "H3Ll0B0T",
        Net: "tcp",
        Addr: "107.170.45.12:3306",
        DBName: "hello-bot",
    }
    return sql.Open("mysql", cfg.FormatDSN())
}

