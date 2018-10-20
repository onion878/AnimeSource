package main

import (
	"./structs"
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

const path = "https://anime1.me"

func main() {
	utils.StartPool()
	//getChapter("https://anime1.me/?cat=333")
	getMenu()
	getAllIndex()
	//utils.SaveOrUpdateIndex("onion", "1-13")
}

func testJSOn() {
	s := `[{"type": "123","file":"213123",label:"435345","default":"56456"}]`
	s = strings.Replace(s, ",label:", `,"label":`, -1)
	var arr []structs.UrlData
	_ = json.Unmarshal([]byte(s), &arr)
	log.Printf("Unmarshaled: %+v\n", arr)
	println(s)
}

func getMenu() {
	c := colly.NewCollector()
	println("获取所有目录")
	c.OnHTML(".entry-content table tbody tr", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Find(".column-1 a").Attr("href")
		name := e.DOM.Find(".column-1 a").Text()
		chapter := e.DOM.Find(".column-2").Text()
		utils.SaveIndex(name, chapter, href, e.Index)
		//index := utils.SaveOrUpdateIndex(name, chapter)
		//getChapter(url + href, index.Id)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(path)
}

func getAllIndex() {
	index := utils.GetAllIndex()
	println("获取详情...")
	for i := 0; i < len(index); i++ {
		data := index[i]
		getChapter(path+data.Url, data.Id)
	}
}

func getChapter(url string, pid string) {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("main", func(e *colly.HTMLElement) {
		s := e.DOM.Find("iframe[src]")
		d := e.DOM.Find(".entry-title a[href]")
		for i := 0; i < s.Length(); i++ {
			src, _ := s.Eq(i).Attr("src")
			name := d.Eq(i).Text()
			getChapterUrl(src, name, pid, i)
		}
	})

	c.OnHTML(".nav-previous a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}

func getChapterUrl(url string, name string, pid string, num int) {
	c := colly.NewCollector(colly.Async(true))
	// Find and visit all links
	c.OnHTML("body script", func(e *colly.HTMLElement) {
		data := e.Text
		start := strings.Index(data, "sources:")
		end := strings.Index(data, ",controls:true")
		if start > 0 && end > 0 {
			s := data[start+8 : end]
			s = strings.Replace(s, ",label:", `,"label":`, -1)
			var arr []structs.UrlData
			_ = json.Unmarshal([]byte(s), &arr)
			var flag = false
			for i := 0; i < len(arr); i++ {
				if arr[i].Default == "true" {
					utils.SaveChapter(name, pid, arr[i].File, num)
					flag = true
				}
			}
			if !flag {
				some := 0
				file := ""
				for i := 0; i < len(arr); i++ {
					s := arr[i].Label
					if len(s) > 0 {
						hd, _ := strconv.Atoi(s[0 : len(s)-1])
						if hd > some {
							file = arr[i].File
						}
						some = hd
					}
				}
				if len(file) > 0 {
					utils.SaveChapter(name, pid, file, num)
				}
			}
		} else {
			start := strings.Index(data, `,file:"`)
			end := strings.Index(data, `",controls:true`)
			if start > 0 && end > 0 {
				file := data[start+7 : end]
				utils.SaveChapter(name, pid, file, num)
			}
		}
	})

	c.OnHTML("video source", func(e *colly.HTMLElement) {
		file := e.Attr("src")
		if len(file) > 0 {
			utils.SaveChapter(name, pid, file, num)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}
