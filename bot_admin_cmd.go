package main

import (
	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func admin_send_all(db *gorm.DB, bot *tgbotapi.BotAPI, text string,
	user *tgbotapi.User, chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("andmin_send_all command")
	if !is_admin {
		return "Permission denied"
	}
	msg, ok := parse_admin_send_all_cmd(text)
	if !ok {
		return "Command arguments error"
	}
	users, err := giveUsers(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "No users found"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	send_all_users(bot, msg, users)
	return "Message send to all users success"
}

func send_all_users(bot *tgbotapi.BotAPI, msg string, users []User) {
	for _, user := range users {
		msg := tgbotapi.NewMessage(user.Chat_id, msg)
		bot.Send(msg)
	}
}

func parse_admin_send_all_cmd(text string) (string, bool) {
	msg := strings.TrimPrefix(text, "/admin_send_all")
	msg = strings.Trim(msg, " ")
	if args := strings.Fields(msg); len(args) < 1 {
		return "", false
	}
	return msg, true
}

func admin_new(db *gorm.DB, text string, user *tgbotapi.User,
	chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("andmin_new command")
	if !is_admin {
		return "Permission denied"
	}
	user_id, ok := parse_admin_cmd("/admin_new", text)
	if !ok {
		return "Command arguments error"
	}
	new_admin, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	if new_admin.UserRole == "admin" {
		return "User with user_id: " + strconv.Itoa(user_id) + " - Is admin"
	}
	new_admin.UserRole = "admin"
	db.Save(&new_admin)
	return "New admin added user_id: " + strconv.Itoa(user_id)
}

func admin_delete(db *gorm.DB, text string, user *tgbotapi.User,
	chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("andmin_new command")
	if !is_admin {
		return "Permission denied"
	}
	user_id, ok := parse_admin_cmd("/admin_delete", text)
	if !ok {
		return "Command arguments error"
	}
	old_adimn, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	if old_adimn.UserRole == "user" {
		return "User with user_id: " + strconv.Itoa(user_id) + " - Not admin"
	}
	old_adimn.UserRole = "user"
	db.Save(&old_adimn)
	return "Delete admin permission user_id: " + strconv.Itoa(user_id)
}

func admin_user_history(db *gorm.DB, text string, user *tgbotapi.User,
	chat_id int64) string {

	is_admin := checkUser(db, user, chat_id)
	log.WithFields(log.Fields{
		"ChatID":   chat_id,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  is_admin,
		"Text":     text,
	}).Info("andmin_user_history command")
	if !is_admin {
		return "Permission denied"
	}
	user_id, ok := parse_admin_cmd("/admin_user_history", text)
	if !ok {
		return "Command arguments errorrrr"
	}
	target, err_fu := giveUserByID(db, user_id)
	if errors.Is(err_fu, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if err_fu != nil {
		log.Error(err_fu)
		return "Sorry something goes wrong, try again"
	}
	ip_req_list, err := giveUserHistory(db, target.User_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "History is empty"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	return userHistory_pp(ip_req_list, target.Name)
}

func parse_admin_cmd(cmd, text string) (int, bool) {
	str_id := strings.TrimPrefix(text, cmd)
	str_id = strings.Trim(str_id, " ")
	if args := strings.Fields(str_id); len(args) != 1 {
		return 0, false
	}
	user_id, err := strconv.Atoi(str_id)
	if err != nil {
		return 0, false
	}
	return user_id, true
}

// Part for admin_all_history
//users, err := giveUsers(db)
//if errors.Is(err, gorm.ErrRecordNotFound) {
//	return "No users found"
//} else if err != nil {
//	log.Error(err)
//	return "Sorry something goes wrong, try again"
//}
//return history_all_users(db, users)
