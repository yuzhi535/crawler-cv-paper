package main

import (
	"crawler/services"
	"fmt"
	"os"
	"strconv"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	sp := selection.New("What do you pick?", []string{"AAAI", "ECCV", "ICCV", "CVPR", "NIPS"})

	sp.PageSize = 5

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	input := textinput.New("Enter Year")
	input.InitialValue = "2022"
	input.Placeholder = "Year cannot be empty"

	year, err := input.RunPrompt()
	if err != nil {
		panic("input valid")
	}
	// year := "2024"

	res, err := strconv.Atoi(year)
	if err != nil {
		panic("Year invalid")
	}

	if choice == "NIPS" {
		services.PipeLine(services.NIPS, res, choice)
	} else if choice == "ECCV" {
		_, err := services.ECCV(res)
		if err != nil {
			panic(err)
		}
		services.PipeLine(services.ECCV, res, choice)
	}

}
