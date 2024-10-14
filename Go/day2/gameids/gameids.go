package gameids

import (
	"bufio"
	"day2/fileops"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CalculateTotal(path string, reader fileops.ReadableFile) (int, error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return -1, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	scanner := bufio.NewScanner(file)

	calibrationValuesTotal := 0

	for scanner.Scan() {
		line := scanner.Text()

		gameId, err := possibleGameIdOrZero(line)
		if err != nil {
			return -1, err
		}

		calibrationValuesTotal += gameId
	}

	return calibrationValuesTotal, nil
}

func possibleGameIdOrZero(line string) (int, error) {
	const maxRed = 12
	const maxGreen = 13
	const maxBlue = 14

	redHighest := 0
	greenHighest := 0
	blueHighest := 0

	revelations := strings.Split(line, ":")
	if len(revelations) < 2 {
		errMsg := fmt.Sprintf("unable to parse line [%s]", line)
		return -1, errors.New(errMsg)
	}

	revelationGroups := strings.Split(revelations[1], ";")
	cubesRegex := regexp.MustCompile(`(\d+)\s(green|blue|red)`)

	for _, group := range revelationGroups {
		matches := cubesRegex.FindAllStringSubmatch(group, -1)
		for _, match := range matches {
			quantity, _ := strconv.Atoi(match[1])
			color := match[2]

			switch color {
			case "red":
				if quantity > redHighest {
					redHighest = quantity
				}
			case "green":
				if quantity > greenHighest {
					greenHighest = quantity
				}
			case "blue":
				if quantity > blueHighest {
					blueHighest = quantity
				}
			}
		}
	}

	gameId, err := extractGameId(line)
	if err != nil {
		errMsg := fmt.Sprintf("unable to extract game id in line [%s]", line)
		return -1, errors.New(errMsg)
	}

	if redHighest <= maxRed && greenHighest <= maxGreen && blueHighest <= maxBlue {
		return gameId, nil
	}

	return 0, nil
}

func extractGameId(line string) (int, error) {
	idRegex := regexp.MustCompile(`Game (\d+):`)
	matches := idRegex.FindStringSubmatch(line)

	if len(matches) < 2 {
		return -1, errors.New("game id not found")
	}

	id, _ := strconv.Atoi(matches[1])

	return id, nil
}
