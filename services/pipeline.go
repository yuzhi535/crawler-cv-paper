package services

import (
	"crawler/models"
	"fmt"
)

func PipeLine(f func(int) ([]models.Paper, error), year int, base string) {
	papers, err := f(year)
	if err != nil {
		fmt.Println(err)
		return
	}
	title := Save2Excel(base, year, papers)
	fmt.Println("\nDone! the excel name is " + title)
}
