package main

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"time"
	"github.com/joeguo/tldextract"
)

type MySQL struct {
	db *sql.DB
}

func NewMySQL(connection string) (mysql *MySQL, err error) {
	mysql = new(MySQL)
	mysql.db, err = sql.Open("mymysql", connection)
	return
}

func (mysql *MySQL) Close() {
	mysql.db.Close()
}

type Domain struct {
	Domain      string
	Da          int
	Available   bool
	Crawled     bool
}

func (mysql *MySQL) StoreDomain(domain Domain) (err error) {
	insert := "insert ignore into domains(domain,da,available,crawled) values(?,?,?,?)"
	_, err = mysql.db.Exec(insert, domain.Domain, domain.Da, domain.Available, domain.Crawled)
	//log.Println(err)
	return
}

func (mysql *MySQL) Crawled(url string) (err error) {
	update := "update urls set status=? where url=?"
	_, err = mysql.db.Exec(update, 1, url)
	return
}

func (mysql *MySQL) Candidate(d Temp) (err error) {
	insert := "insert ignore into qualified(domain,da,tld,category,`time`) values(?,?,?,?,?)"
	result := extract.Extract(d.Domain.Domain)
	if result.Flag == tldextract.Domain {
		r, err := mysql.db.Exec(insert, d.Domain.Domain, d.Da, result.Tld,d.Category, time.Now())
		if err != nil {
			log.Println(err)
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			log.Println(err)
			return err
		}
		inlink := "insert ignore into backlinks(id,url,target) values(?,?,?)"
		_, err = mysql.db.Exec(inlink, id, d.Backlink, d.Target)
		if err != nil {
			log.Println(err)
			return err
		}

	}
	return nil

}



func (mysql *MySQL) GetDomain(d string) ( *Domain, error) {
	q := "select domain,da,available,crawled from domains where domain=?"
	row := mysql.db.QueryRow(q, d)

	domain := &Domain{Domain:d}
	err := row.Scan(&domain.Domain, &domain.Da, &domain.Available, &domain.Crawled)
	if err != nil {
		//doesn't exist
		//log.Println(err)
		return nil, err
	}
	return domain, nil
}

func (mysql *MySQL) ExistDomain(d string) (bool) {
	dm, _ := mysql.GetDomain(d)
	return dm != nil

}



func (mysql *MySQL) Uncrawled(category string) ( []string, error) {
	q := "select url from urls where category=? and status=? limit ?"
	rows, err := mysql.db.Query(q, category, 0, 20)
	if err != nil {
		return nil, err
	}
	ds := make([]string, 0)

	for rows.Next() {
		var d  string
		rows.Scan(&d)
		ds = append(ds, d)
	}
	return ds, nil
}

func (mysql *MySQL) Categories() ( []string, error) {
	q := "select category from urls group by(category)"
	rows, err := mysql.db.Query(q)
	if err != nil {
		return nil, err
	}
	ds := make([]string, 0)

	for rows.Next() {
		var d string
		rows.Scan(&d)
		ds = append(ds, d)
	}
	return ds, nil
}