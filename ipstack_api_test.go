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
