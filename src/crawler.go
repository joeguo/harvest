package main

import (
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"log"
	"bytes"
	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
	"github.com/joeguo/tldextract"

)


func crawler(url string,agent string) (urls []string, err error) {
	internals := make([]string, 0)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("User-Agent", agent)
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
		us := processLinks(doc)
		//log.Println(urls)
		urls = make([]string, 0)

		for _, u := range (us) {
			if !strings.HasPrefix(u, "http") {
				continue
			}
			tldResult := extract.Extract(u)
			if tldResult.Flag == tldextract.Domain {
				d := tldResult.Root + "." + tldResult.Tld
				if strings.Contains(url,d) {
					internals = append(internals, u)
				}else {
					urls = append(urls, u)
				}
			}
		}

	}

	return

}

func processLinks(doc *goquery.Document) (result []string) {
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


