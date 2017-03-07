package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
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
        log.Printf("%+v\n", update)
    }
}
