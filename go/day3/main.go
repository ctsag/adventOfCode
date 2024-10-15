package main

import (
	"adventOfCode/common/fileops"
	"adventOfCode/common/validation"
	"adventOfCode/day3/schematic"
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

	schematicValuesTotal, gearRatiosTotal, err := schematic.CalculateTotals(path, &fileops.FileReader{})
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(2)
		return
	}

	fmt.Printf("The sum of all schematic values is %d\n", schematicValuesTotal)
	fmt.Printf("The sum of all gear ratios is %d\n", gearRatiosTotal)
}
