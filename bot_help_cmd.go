package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
