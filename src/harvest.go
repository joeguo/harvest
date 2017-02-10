package main

import (
	"log"
	"time"
	"runtime"
	"path"
	"strings"
	"fmt"
	"os"
	"os/exec"
	"seomoz"
	//"proxy"
	"errors"
	"sync"
	"github.com/joeguo/dropcatch/name"
	"github.com/joeguo/tldextract"
	"runtime/pprof"

	"io/ioutil"
)

var (
	mysql *MySQL
	//Proxy     *proxy.Proxies
	Moz       *seomoz.Mozs
	extract *tldextract.TLDExtract
	servers = [4]string{"@208.67.222.222", "@208.67.220.220", "@8.8.8.8", "@8.8.4.4"}
	nameAPI *name.NameAPI
	forbidden  []string
)

const (
	StatusNew        = iota
	StatusCrawled
	StatusDomainrx
	DefaultMozFile   = "moz.txt"
	DefaultUserAgent = `Mozilla/5.0 (Windows NT 6.1; rv:15.0) Gecko/20120716 Firefox/15.0a2`
)

type Temp struct {
	Domain
	Category string
	Backlink string
	Target   string
}

func negative(d string) bool {
	for _, word := range (forbidden) {
		//ignore very short forbidden
		if len(word) < 2 {
			continue
		}
		if strings.Contains(d, word) {
			return true
		}
	}
	return false
}

func harvest(thread int , url string,category string, output chan (Temp)) error {
	err := mysql.Crawled(url)
	if err != nil {
		log.Println(err)
	}
	tldResult := extract.Extract(url)
	d := tldResult.Root + "." + tldResult.Tld
	if negative(tldResult.Root) {
		log.Printf("negative:%s\n", d)
		return nil
	}
	output<-Temp{Domain:Domain{Domain:d}, Backlink:url, Target:url}
	if !strings.HasPrefix(url, "http") {
		return nil
	}
	us, err := crawler(url, DefaultUserAgent)
	if err != nil {
		log.Println(err)
		return err
	}
	ds:=make(map[string]bool)
	for _, u := range (us) {
		tldResult = extract.Extract(u)
		d = tldResult.Root+"."+tldResult.Tld
		if negative(tldResult.Root) {
			log.Printf("negative:%s\n", d)
			continue
		}
		if _,ok:=ds[d];!ok{
			//log.Println(d)
			output<-Temp{Domain:Domain{Domain:d},Category:category, Backlink:url, Target:u}
		}
		ds[d]=true

	}
	//err = mysql.StoreDomains(domains)
	log.Printf("%d %s %d", thread, url, len(us))

	return nil
}

func dacheck(thread int, temps []Temp) {
	moz := Moz.Next()
	ds := make([]string, 0)
	for _, t := range (temps) {
		ds = append(ds, t.Domain.Domain)
	}
	avs, err := moz.DomainAuthorities(ds)
	if err != nil {
		avs, err = moz.DomainAuthorities(ds)
	}
	if err != nil {
		return
	}
	for _, t := range (temps) {
		t.Da = int(avs[t.Domain.Domain])
		log.Println("found:",t.Domain.Domain, t.Da)
		err=mysql.Candidate(t)
		if err!=nil{
			log.Println(err)
		}
	}

}

type Hurl struct {
	url string
	category string
}

func harvestThread(thread int , input chan (Hurl), output chan (Temp)) error {

	for {
		url := <-input
		harvest(thread, url.url,url.category, output)
	}
}

func daThread(thread int, input chan (Temp)) {
	temps := make([]Temp, 0)
	for {
		temp := <-input
		temps = append(temps, temp)
		if len(temps) == 2 {
			dacheck(thread, temps)
			temps = make([]Temp, 0)
		}
	}
}

func DNS(domain string) ([]string, error) {
	out, err := exec.Command("dig", servers[0], "+short", "NS", domain).Output()
	if err != nil {
		return []string{}, err
	}
	if len(out) == 0 {
		return []string{}, errors.New("No DNS Data")
	}
	dns := strings.Split(string(out), ".\n")
	if strings.TrimSpace(dns[len(dns)-1]) == "" {
		return dns[0:len(dns)-1], nil
	}
	return dns, nil
}

func checkWithDNSSync(temps *[]Temp) {
	for i, t := range (*temps) {
		dns, _ := DNS(t.Domain.Domain)
		//log.Println(dns)
		if len(dns) == 0 {
			//possible available, we should return it
			(*temps)[i].Domain.Available = true
		}

	}

	log.Println(*temps)
}
func checkWithDNS(temps *[]Temp) {
	var wg sync.WaitGroup
	wg.Add(len(*temps))
	for i, _ := range (*temps) {
		go func(i int) {
			dns, _ := DNS((*temps)[i].Domain.Domain)
			//log.Println(dns)
			if len(dns) == 0 {
				//possible available, we should return it
				(*temps)[i].Available = true
			}
			wg.Done()
		}(i)

	}
	wg.Wait()
}

var tlds = map[string]bool{"me":true, "com":true, "tv":true, "net":true, "org":true, "info":true, "mobi":true}

func check(thread int, temps []Temp, output chan (Temp)) {
	checkWithDNS(&temps)
	for i, t := range (temps) {
		if t.Available {
			r := extract.Extract(t.Domain.Domain)
			if _, o := tlds[r.Tld]; o {
				ok, e := nameAPI.Available(r.Root, r.Tld)
				if e != nil {
					log.Printf("Name.com error %s:%s\n", t.Domain.Domain, e)
					nameAPI.Login()
					ok, _= nameAPI.Available(r.Root, r.Tld)
					temps[i].Available = ok
					continue
				}
				temps[i].Available = ok
			}else {
				temps[i].Available = false
			}


		}
	}

	for _, tt := range (temps) {
		if tt.Available {
			log.Println(tt.Domain.Domain,tt.Available)
			output<-tt
		}else{
			mysql.StoreDomain(tt.Domain)
		}

	}
}

func checkThread(thread int, input chan (Temp), output chan (Temp)) {
	temps := make([]Temp, 0)
	for {
		temp := <-input
		if mysql.ExistDomain(temp.Domain.Domain){
			continue
		}
		temps = append(temps, temp)
		if len(temps) == 10 {
			check(thread, temps, output)
			temps = make([]Temp, 0)
		}
	}
}
func ReadFile(file string) ([]string, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bs), "\n")
	if lines[len(lines) - 1] == "" {
		return lines[0:len(lines) - 1], nil
	}
	return lines, nil
}

func main() {
	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	Moz = seomoz.InitMozs(path.Join(config.RootDir, DefaultMozFile))
	nameAPI = &name.NameAPI{Account:config.NameAcount, Token:config.NameToken}
	ok, err := nameAPI.Login()
	if err != nil || !ok {
		log.Println("can't login name.com")
		return
	}
	forbidden, err = ReadFile(path.Join(config.RootDir, "negative-words.txt"))
	if err != nil {
		log.Printf("Can't read forbidden file: %s, %+v\n", "negative-words.txt", err)
		return
	}
	for i, d := range (forbidden) {
		forbidden[i] = strings.TrimSpace(d)
	}


	mysql, err = NewMySQL(config.MySQLConnection)
	if err != nil {
		log.Printf("mysql error %+v\n", err)
		return
	}
	defer mysql.Close()
	extract,_ = tldextract.New(config.TldCache, false)

	input := make(chan (Hurl), config.ThreadNumber)
	result1 := make(chan (Temp))
	result := make(chan (Temp))

	for i := 0; i < config.ThreadNumber; i++ {
		go harvestThread(i, input, result1)
	}

	for k := 0; k < config.ThreadNumber; k++ {
		go checkThread(k, result1, result)
	}
	for j := 0; j < config.ThreadNumber; j++ {
		go daThread(j, result)
	}


	count := 0
	cs, err := mysql.Categories()
	if err != nil {
		log.Println(err)
	}

	for {
		for _, category := range (cs) {
			urls, err := mysql.Uncrawled(category);
			log.Println(len(urls))
			if err != nil || len(urls) == 0 {
				//close(input)
				//break;
				if err != nil {
					log.Println(err)
				}
				time.Sleep(10 * time.Second)
				continue
			}

			for _, url := range (urls) {

				input<-Hurl{url:url,category:category}
			}
			count+=1
			if (count%500 == 0) && config.Debug {
				WriteRAMProfile(count)

			}
		}
		time.Sleep(3 * time.Second)



	}

}

func WriteRAMProfile(i int) {
	fn := path.Join(config.RootDir, fmt.Sprintf("%d.mprof", i))
	log.Printf("RAM PRofile:%s\n", fn)
	f, err := os.Create(fn)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.WriteHeapProfile(f)
}
