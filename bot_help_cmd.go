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
	msg := help_user() + "\n"
	msg += "---> Admin Help <---\n"
	msg += `"/admin_new [user_id]" - add to user with <user_id> `
	msg += `admin permissions` + "\n"

	msg += `"/admin_delete [user_id]" - take away from user with <user_id> `
	msg += `amdin permissions` + "\n"

	msg += `"/admin_user_history [user_id]" - show all request of user `
	msg += `with <user_id>` + "\n"

	msg += `"/admin_send_all [msg]" - send <msg> to all familiar to `
	msg += `bot users` + "\n"
	return msg
}

func help_user() string {
	msg := "---> User Help <---\n"
	msg += `"/ip [some_ipV4]" - show info about <some_ipV4>` + "\n"
	msg += `"/history" - shows all your requested ips` + "\n"
	return msg
}
