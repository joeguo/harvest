package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/base64"
	"log"
	"crypto/tls"
	"strings"
	"sync"
	"io/ioutil"
	"strconv"
	"io"

)
//http proxy
type Proxy struct {
	User     string
	Password string
	Host     string
	Port     int
}
type Proxies struct {
	proxies []*Proxy
	index   int
	mutex   sync.Mutex
}

func InitProxies(f string) *Proxies {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := string(bs)
	s = strings.TrimSpace(s)
	ps := strings.Split(s, "\n")
	proxies := make([]*Proxy, len(ps))
	for i, p := range (ps) {
		ts := strings.Split(p, ":")
		//log.Println(p)
		port, _ := strconv.Atoi(ts[1])
		proxies[i] = &Proxy{User:ts[2], Password:ts[3], Host:ts[0], Port:port}
	}
	return &Proxies{proxies:proxies}
}

func (this *Proxies) Next() *Proxy {
	this.mutex.Lock()
	defer func() {
		this.index+=1
		this.mutex.Unlock()
	}()
	if this.index >= len(this.proxies) {
		this.index = 0
	}
	return this.proxies[this.index]
}

//don't handle 301 issue
func Get(u string, proxy *Proxy, agent string) (resp *http.Response, err error) {
	var client *http.Client
	if proxy == nil {
		client = http.DefaultClient
	}else {
		var p, parsedU *url.URL
		parsedU, err = url.Parse(u)
		if parsedU.Scheme == "https" {
			p, err = url.Parse(fmt.Sprintf("https://%s:%s@%s:%d", url.QueryEscape(proxy.User), url.QueryEscape(proxy.Password), proxy.Host, proxy.Port))
		}else {
			p, err = url.Parse(fmt.Sprintf("http://%s:%d", proxy.Host, proxy.Port))

		}
		if err != nil {
			log.Println(err)
			return
		}

		transport := &http.Transport{Proxy: http.ProxyURL(p)}
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: transport}
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if proxy != nil && len(proxy.User) > 0 {
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", proxy.User, proxy.Password)))

		req.Header.Set("Proxy-Authorization", basic)
	}
	req.Header.Set("User-Agent", agent)
	//dump,err:=httputil.DumpRequest(req,true)

	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}
	return
}


func Post(u string, proxy *Proxy, bodyType string, body io.Reader, agent string) (resp *http.Response, err error) {
	var client *http.Client
	if proxy == nil {
		client = &http.Client{}
	}else {
		var p, parsedU *url.URL
		parsedU, err = url.Parse(u)
		if parsedU.Scheme == "https" {
			p, err = url.Parse(fmt.Sprintf("https://%s:%s@%s:%d", url.QueryEscape(proxy.User), url.QueryEscape(proxy.Password), proxy.Host, proxy.Port))
		}else {
			p, err = url.Parse(fmt.Sprintf("https://%s:%s@%s:%d", url.QueryEscape(proxy.User), url.QueryEscape(proxy.Password), proxy.Host, proxy.Port))
			//p, err = url.Parse(fmt.Sprintf("http://%s:%d", proxy.Host, proxy.Port))
		}
		if err != nil {
			log.Println(err)
			return
		}
		transport := &http.Transport{Proxy: http.ProxyURL(p)}
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: transport}
	}
	req, _ := http.NewRequest("POST", u, body)
	req.Header.Set("Content-Type", bodyType)
	if proxy != nil && len(proxy.User) > 0 {
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", proxy.User, proxy.Password)))
		req.Header.Add("Proxy-Authorization", basic)
	}
	if len(agent) > 0 {
		req.Header.Set("User-Agent", agent)
	}

	resp, err = client.Do(req)
	return
}





