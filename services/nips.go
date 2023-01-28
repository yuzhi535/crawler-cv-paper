package main

import (
	"crawler/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/xuri/excelize/v2"
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

	f, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(f)

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
	}
	return papers
}

func Save2Excel(base string, year int, papers []models.Paper) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Set value of a cell.
	f.SetColWidth("Sheet1", "B", "C", 60)
	for id, paper := range papers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(id+1), id+1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(id+1), paper.PaperName)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(id+1), paper.URL)
	}

	f.InsertRows("Sheet1", 1, 1)
	f.SetCellValue("Sheet1", "B1", "paper title")
	f.SetCellValue("Sheet1", "C1", "paper url")

	// Save spreadsheet by the given path.
	if err := f.SaveAs(base + strconv.Itoa(year) + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}

func main() {
	papers := NIPS(2022)
	Save2Excel("NIPS", 2022, papers)
}
