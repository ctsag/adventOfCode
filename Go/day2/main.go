package main

import (
	"day2/fileops"
	"day2/gameids"
	"day2/validation"
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

	total, err := gameids.CalculateTotal(path, &fileops.FileReader{})
	if err != nil {
		log.Printf("Error: %s\n", err)
		osExit(2)
		return
	}

	fmt.Printf("The sum of all possible game Ids is %d\n", total)
}
