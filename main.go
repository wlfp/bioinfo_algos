package main

import (
	"fmt"

	"github.com/wlfp/bioinfo_algos/alignment"
)

func main() {

	algorithmChoices := [...]string{"Global alignment"}

	var algorithmChoice int
ChooseAlgorithm:
	fmt.Println("Please select an algorithm, the following choices are available:")
	for index, algorithmName := range algorithmChoices {
		fmt.Printf("%d. %s\n", index, algorithmName)
	}
	fmt.Scanln(&algorithmChoice)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	switch algorithmChoice {
	case 0:
		alignment.Alignment()
	default:
		fmt.Println("That algorithm wasn't recognised, try again?")
		goto ChooseAlgorithm
	}
}
