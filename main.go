package main

import (
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	go botServe(botInit())
	initBackendApi()
	http.ListenAndServe(":3000", nil)
}

func initBackendApi() {
	http.HandleFunc("/API/", homePageAPI)
	http.HandleFunc("/API/get_users", returnAllUsers)
	http.HandleFunc("/API/get_user", returnSingleUser)
	http.HandleFunc("/API/get_history_by_tg", returnSingleUserHistory)
	http.HandleFunc("/API/delete_history_by_tg", deleteHistoryField)
}

func homePageAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePageAPI!")
}
func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Return All Users")
}
func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	user_id := keys[0]
	fmt.Fprintf(w, fmt.Sprintf("Return a User id = %v", user_id))
}
func returnSingleUserHistory(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	user_id := keys[0]
	fmt.Fprintf(w, fmt.Sprintf("Return History for User with id = %v", user_id))
}
func deleteHistoryField(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete history field")
}

func botServe(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.WithFields(log.Fields{
			"UserName": update.Message.From.UserName,
			"Text":     update.Message.Text,
		}).Info("Message from User")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
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
		log.Warn("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	return bot, updates
}
