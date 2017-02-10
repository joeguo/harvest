package crawler

import (
	"fmt"
	"sync"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"log"
	"bytes"
	"time"
	"code.google.com/p/go.net/html"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
	"github.com/joeguo/tldextract"

)

type moreUrls struct{
	depth  int
	urls   []string
}

type empty struct {

}

type Crawler struct {
	Domain   string
	Setting  Setting
	Visited  map[string]bool
	Lock     sync.Mutex
	Done     chan bool
	Result   Result
	more     chan (*moreUrls)
	requests chan (empty)
	extract *tldextract.TLDExtract

}

type Result struct {
	Domains  map[string]string
}

func New(domain string, setting Setting, extract *tldextract.TLDExtract) *Crawler {
	this := &Crawler{Domain:domain, Setting:setting}
	//this.Urls<-fmt.Sprintf("http://www.%s", domain)
	this.Visited = make(map[string]bool)
	this.more = make(chan (*moreUrls), setting.BufferFactor)
	this.extract = extract
	this.requests = make(chan (empty), setting.Requests)
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second*15

	return this
}

func (this *Crawler) Run() {
	depth := 0
	seed := fmt.Sprintf("http://www.%s", this.Domain)
	outstanding := 1
	crawled := 1
	normalizedSeed, err := purell.NormalizeURLString(seed, this.Setting.URLNormalizationFlags)
	if err != nil {
		return
	}
	this.Visited[normalizedSeed] = true
	err = this.crawler(seed, depth)
	if err != nil {
		return
	}
Outer:
	for outstanding > 0 {
		//log.Println(outstanding)
		select {
		case <-this.Done:
			break Outer
		case <-time.After(this.Setting.TimeOut):
			log.Printf("%s timeout", this.Domain)
			break Outer
		case next := <-this.more: {
			outstanding--
			if next.depth > this.Setting.Depth {
				continue
			}
			for _, url := range (next.urls) {

				if crawled > this.Setting.MaxUrls {
					time.Sleep(time.Second * 10)
					return
				}

				normalizedUrl, err := purell.NormalizeURLString(url, this.Setting.URLNormalizationFlags)
				if err != nil {
					continue
				}

				if _, seen := this.Visited[normalizedUrl]; seen {
					continue
				}
				this.Visited[normalizedUrl] = true
				outstanding++
				crawled++

				go this.crawler(url, next.depth)
				<-this.requests
			}
		}
		}


	}

}




func (this *Crawler) crawler(url string, depth int) ( err error) {
	internals := make([]string, 0)
	defer func() {
		this.requests<-empty{}
		this.more<-  &moreUrls{depth + 1, internals}

	}()
	//log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("User-Agent", this.Setting.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return
	}
	//log.Println(contentType)
	bs, err := ioutil.ReadAll(resp.Body)
	if node, e := html.Parse(bytes.NewBuffer(bs)); e != nil {
		log.Println(e)
	} else {
		doc := goquery.NewDocumentFromNode(node)
		doc.Url = resp.Request.URL
		urls := this.processLinks(doc)
		//log.Println(urls)
		if depth == 0 {
			this.Result = Result{ Domains:make(map[string]string)}
		}
		domains := make([]string, 0)

		for _, u := range (urls) {
			if strings.HasPrefix(u, "mail") {
				continue
			}
			tldResult := this.extract.Extract(u)
			if tldResult.Flag == tldextract.Domain {
				d := tldResult.Root + "." + tldResult.Tld
				if d == this.Domain {
					internals = append(internals, u)
				}else {
					domains = append(domains, u)
				}
			}
		}
		if len(domains) > 0 {
			this.Lock.Lock()
			for _, d := range (domains) {
				this.Result.Domains[d] = url
			}
			this.Lock.Unlock()
		}
	}

	//log.Println(body)
	//this.more <- moreUrls{depth + 1, urls}
	return

}

func (this *Crawler) processLinks(doc *goquery.Document) (result []string) {
	urls := doc.Find("a[href]").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Attr("href")
		return val
	})
	for _, s := range urls {
		// If href starts with "#", then it points to this same exact URL, ignore (will fail to parse anyway)
		if len(s) > 0 && !strings.HasPrefix(s, "#") {
			if parsed, e := url.Parse(s); e == nil {
				parsed = doc.Url.ResolveReference(parsed)
				//result = append(result, parsed)
				result = append(result, parsed.String())
			} else {
				//this.logFunc(LogIgnored, "ignore on unparsable policy %s: %s", s, e.Error())
			}
		}
	}
	return
}





