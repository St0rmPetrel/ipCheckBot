package main

import (
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	setDB(db)
	bot, update := botInit()
	go botServe(db, bot, update)
	initBackendApi()
	http.ListenAndServe(":3000", nil)
}

func botInit() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	//bot.Debug = true
	log.WithFields(log.Fields{
		"BotUserName": bot.Self.UserName,
	}).Info("Authorized on account")

	bot.RemoveWebhook()
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(
		os.Getenv("CLOUD_FUNCTION_URL") + "/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.WithFields(log.Fields{
			"LastErrorMessege": info.LastErrorMessage,
		}).Warn("Telgram callback failed")
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	return bot, updates
}
