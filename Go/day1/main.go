package main

import (
	"day1/coordinates"
	"day1/fileops"
	"day1/validation"
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

	total, err := coordinates.CalculateTotal(path, &fileops.FileReader{})
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(2)
		return
	}

	fmt.Printf("The sum of all calibration values is %d\n", total)
}
