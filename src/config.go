package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Debug                       bool
	Key                         string
	RootDir                     string
	MySQLConnection             string
	LmiFolder                   string
	NameSiloAPI                 string
	NameAcount                  string
	NameToken                   string
	TldCache                    string
	ThreadNumber                int
	Depth                       int
	MaxUrls                     int
	Requests                    int
	MinDa                       int
}

const (
	CONFIG_NAME = "config.json"
)

var config *Config

func init() {
	DefaultConfig()
	if Exists(CONFIG_NAME) {
		LoadConfig()
	}else {
		SaveConfig()
	}
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return false
}
func DefaultConfig() {
	config = &Config{Debug:true, Key:"zaq10okm",
		RootDir:"/var/www/crawler/",
		TldCache:"/var/www/crawler/tld.cache",
		MySQLConnection: "crawler/root/apple",
		LmiFolder:"/var/site/share/LMI",
		NameSiloAPI:"namesilo api key",
		NameAcount:"name.com account",
		NameToken:"name.com api token",
		ThreadNumber:16,
		Depth:2,
		MaxUrls:2000,
		Requests:10,
		MinDa:20,
	}
}

func SaveConfig() {
	js, err := json.MarshalIndent(config, "   ", "")
	if err != nil {
		fmt.Println("Marshal json failed:", err)
	}
	ioutil.WriteFile(CONFIG_NAME, js, 0644)
}

func LoadConfig() {
	buff, err := ioutil.ReadFile(CONFIG_NAME)
	if err != nil {
		return
	}
	if err := json.Unmarshal(buff, config); err != nil {
		fmt.Println("Unmarshal json failed:", err)
	}
}

