package main

import "testing"

func TestMoscowIP(t *testing.T) {
	info := NewInfoIP()
	ip := "2.16.103.255" //some Moscow IP
	err := info.fillInfoIP(ip)
	if info.City != "Moscow" || err != nil {
		if err != nil {
			t.Fatalf(`error: %v`, err)
			return
		}
		t.Fatalf(`info.city = "%v", want "Moscow"`, info.City)
	}
}

func TestConnectDB(t *testing.T) {
	_, err := connectToDB()
	if err != nil {
		t.Fatalf(`error: %v`, err)
		return
	}
}

func TestTakeUserFromDB(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fatalf(`error: %v`, err)
		return
	}
	setDB(db)
	var user User
	db.Take(&user)
	if user.Name == "" {
		t.Fatalf(`error: Can't load user from data base`)
		return
	}
}
