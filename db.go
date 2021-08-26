package main

import (
	"encoding/json"
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDB() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable" +
		" TimeZone=Europe/Moscow"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func giveUserByID(db *gorm.DB, user_id int) (*User, error) {
	user := NewUser()
	result := db.Where("user_id = ?", user_id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func giveUsers(db *gorm.DB) ([]User, error) {
	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func giveIpInfoByIP(db *gorm.DB, ip string) (*InfoIP, error) {
	history := NewGlobalHistory()
	result := db.Where("ip = ?", ip).First(&history)
	if result.Error != nil {
		return nil, result.Error
	}
	infoIP := NewInfoIP()
	err := json.Unmarshal([]byte(history.Ip_info), infoIP)
	return infoIP, err
}

func giveUserHistoryByIP(db *gorm.DB,
	ip string, user_id int) (*UserHistory, error) {

	ret := NewUserHistory()
	result := db.Where("user_id = ?", user_id).
		Where("ip = ?", ip).First(&ret)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func giveUserHistoryByID(db *gorm.DB, id int) (*UserHistory, error) {

	ret := NewUserHistory()
	result := db.First(&ret, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func giveUserHistory(db *gorm.DB, user_id int) ([]string, error) {
	var ips []string

	result := db.Table("user_histories").Where("user_id = ?", user_id).
		Select("ip").Find(&ips)
	if result.Error != nil {
		return nil, result.Error
	}
	return ips, nil
}

func giveUserHistoryRet(db *gorm.DB, user_id int) ([]InfoIP, error) {
	ret := []InfoIP{}
	ips, err := giveUserHistory(db, user_id)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		info, err_info := giveIpInfoByIP(db, ip)
		if err_info != nil {
			continue
		}
		ret = append(ret, *info)
	}
	return ret, nil
}

func giveGlobalUserHistory(db *gorm.DB) ([]UserHistory, error) {
	var ret []UserHistory

	result := db.Find(&ret)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func createGlobalHistory(db *gorm.DB, info *InfoIP) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	db.Create(&GlobalHistory{Ip: info.Ip, Ip_info: string(b)})
	return nil
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

func setDB(db *gorm.DB) error {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserHistory{})
	db.AutoMigrate(&GlobalHistory{})
	return nil
}

type User struct {
	gorm.Model
	Name     string
	User_id  int
	Chat_id  int64
	UserRole string
}

func NewUser() *User {
	return &User{}
}

type UserHistory struct {
	gorm.Model
	User_id int
	Ip      string
}

func NewUserHistory() *UserHistory {
	return &UserHistory{}
}

type GlobalHistory struct {
	gorm.Model
	Ip      string
	Ip_info string
}

func NewGlobalHistory() *GlobalHistory {
	return &GlobalHistory{}
}
