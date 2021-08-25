package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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
	ip, ok := parse_ip_cmd(text)
	if !ok {
		return "Command arguments error"
	}
	if ip_info, is_saved := checkIp(db, ip); !is_saved {
		info := NewInfoIP()
		if err := info.fillInfoIP(ip); err != nil {
			log.Error(err)
			return "Sorry something goes wrong, try again"
		}
		if err := createGlobalHistory(db, info); err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"ip": ip,
			}).Info("New IP Info save in Data Base")
		}
		checkUserHistory(db, user, ip)
		return ip_info_pp(info)
	} else {
		log.WithFields(log.Fields{
			"ip": ip,
		}).Info("IP Info exist in Data Base")
		checkUserHistory(db, user, ip)
		return ip_info_pp(ip_info)
	}
}

func parse_ip_cmd(text string) (string, bool) {
	args := strings.Fields(text)
	if len(args) != 2 || !validIP4(args[1]) {
		return "", false
	}
	return args[1], true
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func checkIp(db *gorm.DB, ip string) (*InfoIP, bool) {
	exist_info, err := giveIpInfoByIP(db, ip)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	} else if err != nil {
		log.Error(err)
		return nil, false
	}
	return exist_info, true
}

func ip_info_pp(info *InfoIP) string {
	if info.Country_name == "" {
		ret := "---> ip = " + info.Ip + "\n" +
			"Don't found"
		return ret
	}
	ret := "---> ip = " + info.Ip + "\n" +
		"Country: " + info.Country_name +
		info.Location.Country_flag_emoji + "\n" +
		"City: " + info.City + "\n" +
		"GEO: " + "\n" +
		"--> Latitude:  " + fmt.Sprintf("%v", info.Latitude) + "\n" +
		"--> Longitude: " + fmt.Sprintf("%v", info.Longitude) + "\n"
	return ret
}

func checkUserHistory(db *gorm.DB, user *tgbotapi.User, ip string) {
	_, err := giveUserHistoryByIP(db, ip, user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&UserHistory{User_id: user.ID, Ip: ip})
		log.WithFields(log.Fields{
			"UserID": user.ID,
			"IP":     ip,
		}).Info("Add ip in user history")
		return
	} else if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"UserID": user.ID,
		"IP":     ip,
	}).Info("Ip already exist in user history")
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
