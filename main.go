package main

import (
	"./structs"
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

func main() {
	utils.StartPool()
	//getChapter("https://anime1.me/?cat=333")
	getMenu()
	//utils.SaveOrUpdateIndex("onion", "1-13")
}

func testJSOn()  {
	s := `[{"type": "123","file":"213123",label:"435345","default":"56456"}]`
	s = strings.Replace(s, ",label:", `,"label":`, -1)
	var arr []structs.UrlData
	_ = json.Unmarshal([]byte(s), &arr)
	log.Printf("Unmarshaled: %+v\n", arr)
	println(s)
}

func getMenu() {
	c := colly.NewCollector()
	url := "https://anime1.me"

	c.OnHTML(".entry-content table tbody tr", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Find(".column-1 a").Attr("href")
		name := e.DOM.Find(".column-1 a").Text()
		chapter := e.DOM.Find(".column-2").Text()
		index := utils.SaveOrUpdateIndex(name, chapter)
		getChapter(url + href, index.Id)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
}

func getChapter(url string, pid string) {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("main", func(e *colly.HTMLElement) {
		s := e.DOM.Find("iframe[src]")
		d := e.DOM.Find(".entry-title a[href]")
		for i := 0; i < s.Length(); i++ {
			src,_ :=s.Eq(i).Attr("src")
			name := d.Eq(i).Text()
			getChapterUrl(src, name, pid, i)
		}
	})

	c.OnHTML(".nav-previous a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}

func getChapterUrl(url string, name string, pid string, num int) {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("body script", func(e *colly.HTMLElement) {
		data := e.Text
		start := strings.Index(data, "sources:")
		end := strings.Index(data, ",controls:true")
		if start > 0 && end > 0 {
			s := data[start + 8: end]
			s = strings.Replace(s, ",label:", `,"label":`, -1)
			var arr []structs.UrlData
			_ = json.Unmarshal([]byte(s), &arr)
			for i := 0; i < len(arr); i++ {
				if arr[i].Default == "true" {
					utils.SaveChapter(name, pid, arr[i].File, num)
				}
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}
