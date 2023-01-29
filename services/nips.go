package services

import (
	"crawler/models"
	"log"
	"os"
	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/kinsey40/pbar"
)

func NIPS(year int) []models.Paper {
	soup.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0")
	res, err := soup.Get("https://nips.cc/Conferences/" + strconv.Itoa(year) + "/Schedule?type=Poster")
	if err != nil {
		os.Exit(-1)
	}

	body := soup.HTMLParse(res)
	links := body.FindAllStrict("div", "class", "maincard narrower poster")

	var id = 0
	var papers []models.Paper

	f, err := os.OpenFile("NIPS"+strconv.Itoa(year)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(f)

	p, err := pbar.Pbar(links)
	if err != nil {
		panic(err)
	}
	// Alter pbar settings (e.g. add a description)
	p.SetDescription("Pbar")
	// Initialize just before for-loop
	p.Initialize()

	for _, link := range links {
		title := link.FindStrict("div", "class", "maincardBody").Text()
		href := link.FindStrict("a", "class", "btn btn-default btn-xs href_URL")
		var url string
		if href.Error != nil {
			href = link.FindStrict("a", "title", "Paper PDF")
			url = href.Attrs()["href"]
		} else {
			url = href.Attrs()["href"]
			url = url[:23] + "pdf" + url[28:]
		}
		log.Println(f, "Number:"+strconv.Itoa(id)+" Title: "+title+" url: "+url)
		papers = append(papers, models.Paper{PaperName: title, URL: url})
		id++
		p.Update()
	}
	return papers
}
