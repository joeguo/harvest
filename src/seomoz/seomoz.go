package seomoz

import (
	"time"
	"encoding/base64"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"net/http"
	"log"
	"bytes"
	"seomoz/umcols"
	//"proxy"
	"io/ioutil"
	"sync"
	"strings"

)

const (
	Base = "http://lsapi.seomoz.com/linkscape/%s?%s"
)

type Seomoz struct {
	AccessId    string
	SecretKey   string
	//Proxy *proxy.Proxy
}

func (moz * Seomoz) signature(expires int64) string {
	sign := fmt.Sprintf("%s\n%d", moz.AccessId, expires)
	hash := hmac.New(sha1.New, []byte(moz.SecretKey))
	io.WriteString(hash, sign)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (moz *Seomoz) query(method string, data []byte, params url.Values) (r *http.Response, err error) {
	expires := time.Now().Unix() + 600
	params.Add("AccessID", moz.AccessId)
	params.Add("Expires", strconv.FormatInt(expires, 10))
	params.Add("Signature", moz.signature(expires))
	request := fmt.Sprintf(Base, method, params.Encode())
	buffer := bytes.NewBuffer(data)
	return http.Post(request, "application/x-www-form-urlencoded", buffer)
	//return proxy.Post(moz.Proxy,request, "application/x-www-form-urlencoded", buffer,"")

	//json.NewDecoder(r.Body)
	//decoder.Decode()
}

func (moz *Seomoz) UrlMetrics(urls []string, cols int64) (result []map[string]interface{}, err error) {
	params := url.Values{}

	params.Add("Cols",strconv.FormatInt(cols,10))
	bs, err := json.Marshal(urls)
	if err != nil {
		log.Println(err)
	}
	r, err := moz.query("url-metrics", bs, params)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()
	bs, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(moz.AccessId, string(bs), err)
	}
	err = json.Unmarshal(bs, &result)
	//decoder := json.NewDecoder(r.Body)
	//err = decoder.Decode(&result)
	if err != nil {
		log.Println(moz.AccessId, string(bs), err)
	}

	return
}

func (moz *Seomoz) DomainAuthority(domain string) (uint8, error) {
	result, err := moz.UrlMetrics([]string{domain}, umcols.HttpStatusCode|umcols.DomainAuthority|umcols.Url)
	if err != nil {
		return 0, err
	}
	return uint8(result[0]["pda"].(float64)), err

}

func (moz *Seomoz) DomainAuthorities(domains []string) (map[string]uint8, error) {
	result, err := moz.UrlMetrics(domains, umcols.HttpStatusCode|umcols.DomainAuthority|umcols.Url)
	if err != nil {
		return nil, err
	}
        das:=make(map[string]uint8)
        for i,d:=range(domains){
           das[d]=uint8(result[i]["pda"].(float64))
           //log.Println(d,das[d],result[i])
        }
	return das, err

}

type Mozs struct {
	mozs  []*Seomoz
	index int
	mutex sync.Mutex
}

func InitMozs(f string) *Mozs {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := string(bs)
	ps := strings.Split(s, "\n")
	mozs := make([]*Seomoz, len(ps))
	for i, p := range ps {
		ts := strings.Split(p, ",")
		mozs[i] = &Seomoz{strings.TrimSpace(ts[0]), strings.TrimSpace(ts[1])}
	}
	return &Mozs{mozs: mozs}
}

func (this *Mozs) Next() *Seomoz {
	this.mutex.Lock()
	defer func() {
		this.index += 1
		this.mutex.Unlock()
	}()
	if this.index >= len(this.mozs) {
		this.index = 0
	}
	return this.mozs[this.index]
}
