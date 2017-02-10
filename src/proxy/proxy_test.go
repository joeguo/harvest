package proxy

import (
	"testing"
	"log"
	"io/ioutil"
	"net/url"
	"strings"

)

func TestGet(t *testing.T) {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	agent   := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.43 Safari/537.31"
	proxy := &Proxy{User:"a", Password:"aa#", Host:"199.193.255.39", Port:43543}
	resp,err:=Get("http://httpbin.org/headers", proxy, agent)
	if err!=nil{
       log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

	}
	r := string(body)
	log.Println(r)
}
func TestPost(t *testing.T) {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	agent   := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.43 Safari/537.31"
	proxy := &Proxy{User:"aa", Password:"aaa#", Host:"199.193.255.39", Port:43543}
	form:=url.Values{"key": {"Value"}, "id": {"123"}}
	resp,err:=Post("http://httpbin.org/post",proxy, "application/x-www-form-urlencoded",strings.NewReader(form.Encode()), agent)
	if err!=nil{
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

	}
	r := string(body)
	log.Println(r)
}