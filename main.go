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

	res, err := strconv.Atoi(year)
	if err != nil {
		panic("Year invalid")
	}

	if choice == "NIPS" {
		services.PipeLine(services.NIPS, res, choice)
	} else if choice == "ECCV" {
		services.PipeLine(services.ECCV, res, choice)
	}

}
