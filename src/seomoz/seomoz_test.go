package seomoz

import (
	"testing"
	"seomoz/umcols"
	"log"
	"proxy"
)

func init(){
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
var (

	Proxy = proxy.InitProxies("/home/joe/dev/go/domainrankhistory/src/history/proxy.txt")
	seo = &Seomoz{"member-78b2986584", "c286f6958f2198d486bc09fa452625d2"}
)

func TestSignature(t *testing.T) {
	signature := seo.signature(1369632518)
	if signature != "wFVu0rxcs/PdIwQw62v9fXMTTP4=" {
		t.FailNow()
	}
}

func TestUrlMetrics(t *testing.T) {
	links := []string{"joomlatags.org", "hostye.com", "rankquote.com"}
	result, err := seo.UrlMetrics(links, umcols.HttpStatusCode | umcols.DomainAuthority | umcols.Url)
	log.Println(result, err)
}

func aTestDomainAuthority(t *testing.T) {
	domain := "google.com"
	result, err := seo.DomainAuthority(domain)
	log.Println(result, err)
}
