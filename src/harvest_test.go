package main

import (
	"testing"
    "log"
	"path"
	"seomoz"
	"registrar/namesilo"
	"github.com/joeguo/tldextract"
	"time"

)

func TestCheck(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	Moz = seomoz.InitMozs(path.Join(config.RootDir, DefaultMozFile))
	namesiloAPI = &namesilo.NamesiloAPI{Key:config.NameSiloAPI}


	mysql, _ = NewMySQL(config.MySQLConnection)

	extract = tldextract.New(config.TldCache, false)
	temps:=make([]Temp,0)
	temps=append(temps,Temp{Domain:Domain{Domain:"usershare.net",Da:25}, Backlink:"hostye.com", Target:"hostye123.com"})
	check(1,temps)
	time.Sleep(time.Minute)

}


