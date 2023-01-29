package services

import (
	"crawler/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/kinsey40/pbar"
	"github.com/xuri/excelize/v2"
)

func ECCV(year int) []models.Paper {
	soup.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0")
	resp, err := soup.Get("https://eccv" + strconv.Itoa(year) + ".ecva.net/program/accepted-papers/")
	if err != nil {
		os.Exit(1)
	}
	body := soup.HTMLParse(resp)

	// Get all the rows
	rows := body.FindStrict("div", "class", "entry-content")
	if rows.Error != nil {
		// If there are no papers, exit
		os.Exit(-1)
	}

	f, err := os.OpenFile("ECCV"+strconv.Itoa(year)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(f)

	var papers []models.Paper

	trs := rows.FindAll("tr")

	p, err := pbar.Pbar(trs)
	if err != nil {
		panic(err)
	}
	// Alter pbar settings (e.g. add a description)
	p.SetDescription("Pbar")
	// Initialize just before for-loop
	p.Initialize()

	// Get all the rows
	for _, paper := range trs {
		if paper.Attrs()["style"] == "background-color: #eee" {
			continue
		}
		tds := paper.FindAll("td")
		// Get the ID

		id := tds[0].Text()
		// Get the title
		title := tds[1].Text()
		// Get the pdf link
		link := "https://www.ecva.net/papers/eccv_2022/papers_ECCV/html/" + id + "_ECCV_2022_paper.php"
		pdf := FetchPDF(link)
		// Print the paper
		log.Println("id: " + id + " title: " + title + " url: " + pdf)
		papers = append(papers, models.Paper{PaperName: title, URL: pdf})
		p.Update()
	}
	return papers
}

func FetchPDF(url string) (pdf string) {
	// soup.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0")
	resp, err := soup.Get(url)
	if err != nil {
		os.Exit(1)
	}
	body := soup.HTMLParse(resp)

	// Get all the rows
	rows := body.FindStrict("div", "id", "content")
	if rows.Error != nil {
		// If there are no papers, exit
		os.Exit(-1)
	}

	// Get all the rows
	link := rows.Find("a").Attrs()["href"]
	pdf = url[:50] + "papers/" + link[48:]
	return pdf
}

// Path: main.go
func main() {
	year := 2022
	papers := ECCV(year)
	base := "ECCV"

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
	title := base + strconv.Itoa(year) + ".xlsx"
	if err := f.SaveAs(base + strconv.Itoa(year) + ".xlsx"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done! the excel name is " + title)
}
