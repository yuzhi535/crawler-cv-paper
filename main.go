package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	base := promptui.Select{
		Label: "Select Conference",
		Items: []string{"AAAI", "CVPR", "ECCV", "ICCV", "NIPS",
			"ICLR"},
	}

	_, result, err := base.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

}
