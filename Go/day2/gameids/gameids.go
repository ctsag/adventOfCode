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

func CalculateTotals(path string, reader fileops.ReadableFile) (idTotal int, minCubesTotal int, err error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return -1, -1, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		minCubes, err := minimumPossibleCubes(line)
		if err != nil {
			return -1, -1, err
		}

		minCubesTotal += minCubes

		gameId, err := possibleGameIdOrZero(line)
		if err != nil {
			return -1, -1, err
		}

		idTotal += gameId
	}

	return idTotal, minCubesTotal, nil
}

func possibleGameIdOrZero(line string) (int, error) {
	const maxRed = 12
	const maxGreen = 13
	const maxBlue = 14

	gameId, err := extractGameId(line)
	if err != nil {
		errMsg := fmt.Sprintf("unable to extract game id in line [%s]", line)
		return -1, errors.New(errMsg)
	}

	highestPerColor, err := extractMaxColorValues(line)
	if err != nil {
		errMsg := fmt.Sprintf("unable to extract color values in line [%s]", line)
		return -1, errors.New(errMsg)
	}

	if highestPerColor["red"] <= maxRed && highestPerColor["green"] <= maxGreen && highestPerColor["blue"] <= maxBlue {
		return gameId, nil
	}

	return 0, nil
}

func extractMaxColorValues(line string) (map[string]int, error) {
	highestPerColor := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}

	revelations := strings.Split(line, ":")
	if len(revelations) < 2 {
		errMsg := fmt.Sprintf("unable to parse line [%s]", line)
		return nil, errors.New(errMsg)
	}

	revelationGroups := strings.Split(revelations[1], ";")

	cubesRegex := regexp.MustCompile(`(\d+)\s(green|blue|red)`)

	for _, group := range revelationGroups {
		matches := cubesRegex.FindAllStringSubmatch(group, -1)
		for _, match := range matches {
			quantity, _ := strconv.Atoi(match[1])
			color := match[2]

			if quantity > highestPerColor[color] {
				highestPerColor[color] = quantity
			}
		}
	}

	return highestPerColor, nil
}

func extractGameId(line string) (int, error) {
	idRegex := regexp.MustCompile(`Game (\d+)`)
	matches := idRegex.FindStringSubmatch(line)

	if len(matches) < 2 {
		return -1, errors.New("game id not found")
	}

	id, _ := strconv.Atoi(matches[1])

	return id, nil
}

func minimumPossibleCubes(line string) (int, error) {
	highestPerColor, err := extractMaxColorValues(line)
	if err != nil {
		errMsg := fmt.Sprintf("unable to extract color values in line [%s]", line)
		return -1, errors.New(errMsg)
	}

	minCubes := highestPerColor["red"] * highestPerColor["green"] * highestPerColor["blue"]

	return minCubes, nil
}
