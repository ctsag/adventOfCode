package main

import (
	"adventOfCode/common/fileops"
	"adventOfCode/common/validation"
	"adventOfCode/day4/scratchcards"
	"fmt"
	"log"
	"os"
)

var osExit = os.Exit

func main() {
	log.SetFlags(0)

	path, err := validation.ExtractSingleArgIgnoringOthers(os.Args, 2)
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(1)
		return
	}

	scratchcardsTotal, scratchcardCount, err := scratchcards.CalculateTotals(path, &fileops.FileReader{})
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(2)
		return
	}

	fmt.Printf("The sum of all scratchcards is %d\n", scratchcardsTotal)
	fmt.Printf("The count of all bonus scratchcards is %d\n", scratchcardCount)
}
