package services

import (
	"crawler/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kinsey40/pbar"
)

func getCurrentYear() int {
	return time.Now().Year()
}

func extractEccvId(input string) string {
	re := regexp.MustCompile(`(\d+)_ECCV`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func ECCV(year int) ([]models.Paper, error) {
	url := "https://www.ecva.net/papers.php"
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

	f, err := os.OpenFile("ECCV"+strconv.Itoa(year)+".log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Get all the rows
	rows := contents.Find("button.accordion")
	if rows.Length() == 0 {
		log.Println("there are no papers!")
		return nil, fmt.Errorf("no papers found")
	}

	currentYear := getCurrentYear()
	idx := (currentYear - year) / 2

	if idx >= rows.Length() || idx < 0 || year < 2018 {
		log.Println("not a valid year!")
		return nil, fmt.Errorf("not a valid year")
	}

	row_of_papers := rows.Eq(idx).Next()
	raw_papers := row_of_papers.Find("dt.ptitle")

	if raw_papers.Length() == 0 {
		log.Printf("there are no papers in %d!\n", year)
		return nil, fmt.Errorf("no papers found in %d", year)
	}

	var papers []models.Paper
	p, err := pbar.Pbar(raw_papers.Nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to create pbar: %v", err)
	}
	p.SetDescription("Processing ECCV papers")
	p.Initialize()
	raw_papers.Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, _ := s.Find("a").Attr("href")
		id := extractEccvId(link)
		pdf := "https://www.ecva.net/" + link
		log.Printf("Id: %d, Title: %s, pdf: %s \n", id, title, pdf)
		paper := models.Paper{PaperName: title, URL: pdf}
		papers = append(papers, paper)
		p.Update()
	})

	return papers, nil
}
