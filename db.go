package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDB() (*gorm.DB, error) {
	dsn := "host=localhost" +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable" +
		" TimeZone=Europe/Moscow"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func giveUserByID(db *gorm.DB, user_id uint) (*User, error) {
	user := NewUser()
	result := db.Where("user_id = ?", user_id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
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

func giveUserHistory(db *gorm.DB, user_id uint) ([]string, error) {
	var ips []string

	result := db.Table("UserHistory").Where("user_id = ?", user_id).
		Select("ip").Find(&ips)
	if result.Error != nil {
		return nil, result.Error
	}
	return ips, nil
}

func giveGlobalHistory(db *gorm.DB) ([]string, error) {
	var ips []string

	result := db.Table("GlobalHistory").Select("ip").Find(&ips)
	if result.Error != nil {
		return nil, result.Error
	}
	return ips, nil
}

func createGlobalHistory(db *gorm.DB, info InfoIP) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	db.Create(&GlobalHistory{Ip: info.Ip, Ip_info: string(b)})
	return nil
}

func setDB(db *gorm.DB) error {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserHistory{})
	db.AutoMigrate(&GlobalHistory{})
	admin_user_id, _ := strconv.Atoi(os.Getenv("ADMIN_USER_ID"))
	_, err := giveUserByID(db, uint(admin_user_id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&User{Name: "Creator", User_id: uint(admin_user_id),
			UserRole: "admin"})
	} else {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name     string
	User_id  uint
	UserRole string
}

func NewUser() *User {
	return &User{}
}

type UserHistory struct {
	gorm.Model
	User_id uint
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
