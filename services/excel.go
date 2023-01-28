package services

import (
	"crawler/models"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func Save2Excel(base string, year int, papers []models.Paper) string {
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
	return title
}
