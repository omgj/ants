package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

const (
	abc      = `abcdefghijklmnopqrstuvwxyz`
	api      = `https://www.etymonline.com/search?page=%d&q=%s`
	sqlwords = `CREATE TABLE IF NOT EXISTS ety (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		word VARCHAR(35),
		def TEXT,
		url VARCHAR(50)
	)`
)

var db *sql.DB

// ascicolors
const (
	reset   = `\u001b[0m`
	magenta = `\u001b[35m`
	blue    = `\u001b[34m`
	green   = `\u001b[32m`
	yellow  = `\u001b[33mp`
	cyan    = `\u001b[36m`
	red     = `\u001b[31m`
)

func hand(w http.ResponseWriter, r *http.Request) {
	home := false
	a := strings.Split(r.URL.String(), "/")
	var q string
	if len(a) == 2 {
		q = a[1]
	}
	if q == `` {
		home = true
	}

	if home {
		t, e := template.ParseFiles("index.html")
		if e != nil {
			log.Println(e)
			return
		}
		e = t.Execute(w, nil)
		if e != nil {
			log.Println(e)
			return
		}
		return
	}

	fmt.Println("Searching for word in cache: ", q)

}

func main() {
	var e error
	db, e = sql.Open("mysql", "root:@/words")
	if e != nil {
		log.Println(e)
		return
	}
	_, e = db.Exec(sqlwords)
	if e != nil {
		log.Println(e)
		return
	}

	c := colly.NewCollector()

	var nothing bool
	c.OnHTML("div.word--C9UPa", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a", "href")
		def := e.ChildText("section.word__defination--2q7ZH")
		i := strings.Split(strings.Split(url, "#")[0], "/")
		word := i[len(i)-1]
		_, ee := db.Exec(`insert into ety (word, def, url) values (?,?,?)`, word, def, url)
		if ee != nil {
			log.Println(ee)
			return
		}
		g := len(def)
		if g > 10 {
			def = def[:10]
		}
		fmt.Printf("%s\t%s\n", word, def)
		nothing = false
	})

	for _, b := range abc {
		for i := 1; i < 1000; i++ {
			nothing = true
			u := fmt.Sprintf(api, i, string(b))
			fmt.Println("visiting ", u)
			c.Visit(u)
			if nothing {
				i = 1000
			}
			time.Sleep(time.Second / 8)
		}
	}

}

func terminalprint() {
	for i := 0; i < 17; i++ {
		for j := 0; j < 17; j++ {
			a := fmt.Sprintf("%d", i*16+j)
			fmt.Printf(`\u001b[38;5;%sm A %s`, a, reset)
		}
	}
}
