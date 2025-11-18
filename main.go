package main

import (
	"fmt"

	"github.com/wlfp/bioinfo_algos/alignment"
	"github.com/wlfp/bioinfo_algos/phylogeny"
)

func main() {

	algorithmChoices := [...]string{"Global alignment", "UPGMA"}

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
	case 1:
		phylogeny.UPGMA()
	default:
		fmt.Println("That algorithm wasn't recognised, try again?")
		goto ChooseAlgorithm
	}
}
