package main

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func botServe(db *gorm.DB, bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel) {

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.WithFields(log.Fields{
			"UserName": update.Message.From.UserName,
			"Text":     update.Message.Text,
		}).Info("Message from User")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				msg.Text = help(db, update.Message.Text,
					update.Message.From, update.Message.Chat.ID)
			case "ip":
				msg.Text = ip(db, update.Message.Text,
					update.Message.From, update.Message.Chat.ID)
			case "history":
				msg.Text = history(db, update.Message.Text,
					update.Message.From, update.Message.Chat.ID)
			default:
				msg.Text = unknown(db, update.Message.Text,
					update.Message.From, update.Message.Chat.ID)
			}
		} else {
			msg.Text = unknown(db, update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		}
		bot.Send(msg)
	}
}

func checkUser(db *gorm.DB, user *tgbotapi.User, chat_id int64) bool {
	exist_user, err := giveUserByID(db, user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&User{Name: user.UserName, User_id: user.ID,
			Chat_id: chat_id, UserRole: "user"})
		log.WithFields(log.Fields{
			"ChatID":   chat_id,
			"UserID":   user.ID,
			"UserName": user.UserName,
		}).Info("Add a new user")
		return false
	} else if err != nil {
		log.Error(err)
		return false
	}
	return exist_user.UserRole == "admin"
}

func ip(db *gorm.DB, text string, user *tgbotapi.User, chat_id int64) string {
	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("ip comand")
	return "ip comand"
}

func history(db *gorm.DB, text string, user *tgbotapi.User,
	chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("history comand")
	return "history comand"
}

func help(db *gorm.DB, text string, user *tgbotapi.User, chat_id int64) string {
	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("help command")
	if is_admin {
		return help_admin()
	}
	return help_user()
}

func help_admin() string {
	msg := "Help admin msg"
	return msg
}

func help_user() string {
	msg := "Help user msg"
	return msg
}

func unknown(db *gorm.DB, text string, user *tgbotapi.User,
	chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Warn("Unknow request")
	return "I don't know that command, try \"/help\""
}
