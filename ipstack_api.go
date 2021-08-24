package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	log.Info("Receive data from ipstack api")
	r, err := myClient.Get(url)
	if err != nil {
		log.Warn("No response from request")
		return err
	}
	defer r.Body.Close()
	log.Info("Decode recive data to InfoIP struct")
	return json.NewDecoder(r.Body).Decode(target)
}

func NewInfoIP() *InfoIP {
	log.Info("Creating new InfoIP")
	return &InfoIP{}
}

func (info *InfoIP) fillInfoIP(ip string) error {
	var url string = "http://api.ipstack.com/" + ip +
		"?access_key=" + os.Getenv("IPSTACK_ACCESS_KEY")

	return getJson(url, info)
}

type InfoIP struct {
	City           string   `json:"city"`
	Continent_code string   `json:"continent_code"`
	Continent_name string   `json:"continent_name"`
	Country_code   string   `json:"country_code"`
	Country_name   string   `json:"country_name"`
	Ip             string   `json:"ip"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	Location       Location `json:"location"`
	Region_code    string   `json:"region_code"`
	Region_name    string   `json:"region_name"`
	Ip_type        string   `json:"type"`
	Zip            string   `json:"zip"`
}

type Location struct {
	Calling_code               string     `json:"calling_code"`
	Capital                    string     `json:"capital"`
	Country_flag               string     `json:"country_flag"`
	Country_flag_emoji         string     `json:"country_flag_emoji"`
	Country_flag_emoji_unicode string     `json:"country_flag_emoji_unicode"`
	Geoname_id                 int        `json:"geoname_id"`
	Is_eu                      bool       `json:"is_eu"`
	Languages                  []Language `json:"languages"`
}

type Language struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Native string `json:"native"`
}
