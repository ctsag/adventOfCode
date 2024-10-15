package main

import (
	"adventOfCode/common/fileops"
	"adventOfCode/common/validation"
	"adventOfCode/day2/gameids"
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

	idsTotal, minCubesTotal, err := gameids.CalculateTotals(path, &fileops.FileReader{})
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(2)
		return
	}

	fmt.Printf("The sum of all possible game ids is %d\n", idsTotal)
	fmt.Printf("The sum of the minimum posisble cubes is %d\n", minCubesTotal)
}
