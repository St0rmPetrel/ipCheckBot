package main

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
	if ok := parse_history_cmd(text); !ok {
		return "Command arguments error"
	}
	ip_req_list, err := giveUserHistory(db, user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "History is empty"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	return userHistory_pp(ip_req_list, user.UserName)
}

func userHistory_pp(ip_list []string, user_name string) string {
	str := "History of user: " + user_name + "\n"
	if len(ip_list) < 1 {
		str += "Empty"
		return str
	}
	for _, ip := range ip_list {
		str += ip + "\n"
	}
	return str + "\n"
}

func parse_history_cmd(text string) bool {
	args := strings.Fields(text)
	if len(args) != 1 {
		return false
	}
	return true
}

func history_all_users(db *gorm.DB, users []User) string {
	ret := ""
	for _, user := range users {
		ip_req_list, err := giveUserHistory(db, user.User_id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err)
			continue
		}
		ret += userHistory_pp(ip_req_list, user.Name)
	}
	return ret
}
