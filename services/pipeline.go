package services

import (
	"crawler/models"
	"fmt"
)

func PipeLine(f func(int) []models.Paper, year int, base string) {
	papers := f(year)
	title := Save2Excel(base, year, papers)
	fmt.Println("Done! the excel name is " + title)
}
