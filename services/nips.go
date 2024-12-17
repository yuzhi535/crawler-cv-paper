package services

import (
	"crawler/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/kinsey40/pbar"
)

func NIPS(year int) ([]models.Paper, error) {
	url := "https://nips.cc/Conferences/" + strconv.Itoa(year) + "/Schedule?type=Poster"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	contents, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	f, err := os.OpenFile("NIPS"+strconv.Itoa(year)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	links := contents.Find("div.maincard.narrower.poster")
	if links.Length() == 0 {
		log.Println("no posters found")
		return nil, fmt.Errorf("no posters found")
	}

	p, err := pbar.Pbar(links.Nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to create pbar: %v", err)
	}
	p.SetDescription("Pbar")
	p.Initialize()

	var id = 0
	var papers []models.Paper

	links.Each(func(i int, s *goquery.Selection) {
		title := s.Find("div.maincardBody").Text()
		href := s.Find("a.btn.btn-default.btn-xs.href_URL")
		var url string
		if href.Length() == 0 {
			href = s.Find("a[title='Paper PDF']")
			url = href.AttrOr("href", "")
		} else {
			url = href.AttrOr("href", "")
			url = url[:23] + "pdf" + url[28:]
		}
		log.Printf("Number: %d Title: %s url: %s\n", id, title, url)
		papers = append(papers, models.Paper{PaperName: title, URL: url})
		id++
		p.Update()
	})

	return papers, nil
}
